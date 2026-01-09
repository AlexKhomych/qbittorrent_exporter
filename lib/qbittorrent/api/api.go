package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"qbittorrent_exporter/types"
	"strings"
)

const (
	apiV2 = "/api/v2"

	authLogin = apiV2 + "/auth/login"

	torrentsInfo = apiV2 + "/torrents/info"
	transferInfo = apiV2 + "/transfer/info"
	appVersion   = apiV2 + "/app/version"

	headerContentType      = "Content-Type"
	headerReferer          = "Referer"
	contentTypeFormEncoded = "application/x-www-form-urlencoded"
	contentTypeJSON        = "application/json"
	contentTypePlain       = "text/plain; charset=UTF-8"
)

type QBittorrentAPI struct {
	baseURL   string
	sidCookie *http.Cookie
	client    *http.Client
}

type QBittorrentAPIOpts struct {
	BaseURL     string
	Credentials *QBittorrentCredentials
	HttpClient  *http.Client
}

type QBittorrentCredentials struct {
	Username string
	Password string
}

func NewQBittorrentAPI(o *QBittorrentAPIOpts) (*QBittorrentAPI, error) {
	api := &QBittorrentAPI{
		baseURL: o.BaseURL,
		client:  o.HttpClient,
	}

	credentials := url.Values{
		"username": {o.Credentials.Username},
		"password": {o.Credentials.Password},
	}
	o.Credentials = &QBittorrentCredentials{}

	if err := api.Login(credentials); err != nil {
		return nil, err
	}

	return api, nil
}

func (api *QBittorrentAPI) Login(credentials url.Values) error {
	var sidCookie *http.Cookie

	loginURL := api.baseURL + authLogin
	if err := ValidateURL(loginURL); err != nil {
		return fmt.Errorf("invalid login URL: %w", err)
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(credentials.Encode()))
	if err != nil {
		return fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set(headerContentType, contentTypeFormEncoded)
	req.Header.Set(headerReferer, api.baseURL)

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "SID" {
			sidCookie = cookie
			break
		}
	}

	if sidCookie == nil {
		return fmt.Errorf("SID cookie not found in login response")
	}
	api.sidCookie = sidCookie
	return nil
}

func (api *QBittorrentAPI) doAuthenticatedGet(endpoint, contentType string) ([]byte, error) {
	url := api.baseURL + endpoint
	if err := ValidateURL(url); err != nil {
		return nil, fmt.Errorf("invalid URL for %s: %w", endpoint, err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for %s: %w", endpoint, err)
	}

	req.AddCookie(api.sidCookie)
	req.Header.Set(headerContentType, contentType)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body for %s: %w", endpoint, err)
	}

	return body, nil
}

func (api *QBittorrentAPI) TorrentsInfo() ([]types.Torrent, error) {
	var torrents []types.Torrent

	body, err := api.doAuthenticatedGet(torrentsInfo, contentTypeJSON)
	if err != nil {
		return torrents, err
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&torrents); err != nil {
		return torrents, fmt.Errorf("decode torrents info: %w", err)
	}

	return torrents, nil
}

func (api *QBittorrentAPI) TransferInfo() (types.Transfer, error) {
	var transfer types.Transfer

	body, err := api.doAuthenticatedGet(transferInfo, contentTypeJSON)
	if err != nil {
		return transfer, err
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&transfer); err != nil {
		return transfer, fmt.Errorf("decode transfer info: %w", err)
	}

	return transfer, nil
}

func (api *QBittorrentAPI) AppVersion() (string, error) {
	body, err := api.doAuthenticatedGet(appVersion, contentTypePlain)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ValidateURL(input string) error {
	_, err := url.ParseRequestURI(input)
	return err
}
