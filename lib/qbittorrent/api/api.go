package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/types"
	"strings"
)

const (
	QBT_APIV2 string = "/api/v2"

	QBT_APIV2_AUTH_LOGIN    = QBT_APIV2 + "/auth/login"
	QBT_APIV2_TORRENTS_INFO = QBT_APIV2 + "/torrents/info"
	QBT_APIV2_TRANSFER_INFO = QBT_APIV2 + "/transfer/info"
	QBT_APIV2_APP_VERSION   = QBT_APIV2 + "/app/version"
)

type QBittorrentAPI struct {
	baseUrl   string
	sidCookie *http.Cookie
}

type QBittorrentAPIOpts struct {
	BaseURL     string
	Credentials *QBittorrentCredentials
}

type QBittorrentCredentials struct {
	Username string
	Password string
}

func NewQBittorrentAPI(o *QBittorrentAPIOpts) (*QBittorrentAPI, error) {
	api := &QBittorrentAPI{
		baseUrl: o.BaseURL,
	}

	credentials := url.Values{
		"username": {o.Credentials.Username},
		"password": {o.Credentials.Password},
	}
	o.Credentials = &QBittorrentCredentials{}

	if err := api.Login(credentials); err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("Failed to login")
	}

	return api, nil
}

func (api *QBittorrentAPI) Login(credentials url.Values) error {
	var sidCookie *http.Cookie

	loginUrl := api.baseUrl + QBT_APIV2_AUTH_LOGIN
	if err := ValidateURL(loginUrl); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to validate auth/login URL")
	}

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(credentials.Encode()))
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to generate auth/login request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", api.baseUrl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to POST auth/login")
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
		return fmt.Errorf("SID Cookie is empty")
	}
	api.sidCookie = sidCookie
	return nil
}

func (api *QBittorrentAPI) TorrentsInfo() ([]types.Torrent, error) {
	var torrents []types.Torrent
	infoUrl := api.baseUrl + QBT_APIV2_TORRENTS_INFO
	if err := ValidateURL(infoUrl); err != nil {
		log.Error(err.Error())
		return torrents, fmt.Errorf("Failed to validate torrents/info URL")
	}

	req, err := http.NewRequest("GET", infoUrl, nil)
	if err != nil {
		log.Error(err.Error())
		return torrents, fmt.Errorf("Failed to generate torrents/info request")
	}
	req.AddCookie(api.sidCookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return torrents, fmt.Errorf("Failed to GET torrents/info")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return torrents, fmt.Errorf("Failed to read torrents/info body")
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&torrents); err != nil {
		log.Error(err.Error())
		return torrents, fmt.Errorf("Failed to decode torrents/info body")
	}

	return torrents, nil
}

func (api *QBittorrentAPI) TransferInfo() (types.Transfer, error) {
	var transfer types.Transfer
	infoUrl := api.baseUrl + QBT_APIV2_TRANSFER_INFO
	if err := ValidateURL(infoUrl); err != nil {
		log.Error(err.Error())
		return transfer, fmt.Errorf("Failed to validate transfer/info URL")
	}

	req, err := http.NewRequest("GET", infoUrl, nil)
	if err != nil {
		log.Error(err.Error())
		return transfer, fmt.Errorf("Failed to generate transfer/info request")
	}
	req.AddCookie(api.sidCookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return transfer, fmt.Errorf("Failed to GET transfer/info")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return transfer, fmt.Errorf("Failed to read transfer/info body")
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&transfer); err != nil {
		log.Error(err.Error())
		return transfer, fmt.Errorf("Failed to decode transfer/info body")
	}

	return transfer, nil
}

func (api *QBittorrentAPI) AppVersion() (string, error) {
	var version string
	versionUrl := api.baseUrl + QBT_APIV2_APP_VERSION
	if err := ValidateURL(versionUrl); err != nil {
		log.Error(err.Error())
		return version, fmt.Errorf("Failed to validate app/version URL")
	}

	req, err := http.NewRequest("GET", versionUrl, nil)
	if err != nil {
		log.Error(err.Error())
		return version, fmt.Errorf("Failed to generate app/version request")
	}
	req.AddCookie(api.sidCookie)
	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return version, fmt.Errorf("Failed to GET app/version")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return version, fmt.Errorf("Failed to read app/version body")
	}
	version = string(body)

	return version, nil
}

func ValidateURL(input string) error {
	_, err := url.ParseRequestURI(input)
	return err
}
