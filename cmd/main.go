package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"qbittorrent_exporter/config"
	"qbittorrent_exporter/feature"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/lib/qbittorrent/api"
	"qbittorrent_exporter/lib/scheduler"
	"qbittorrent_exporter/metrics"
	"qbittorrent_exporter/state"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version = "1.0.2"

const (
	torrentUpdateInterval  = 30 * time.Second
	transferUpdateInterval = 30 * time.Second
	versionCheckInterval   = 10 * time.Minute
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "\nVersion: %v\n\nUsage: %v [ Options... ]\n\nAvailable Options:\n",
			version, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	var (
		logLevel      string
		logFormat     string
		configPath    string
		metricsPrefix string
		useFeatures   = map[feature.FeatureFlag]bool{
			feature.TRANSIENT_STATE: false,
		}
	)

	flag.StringVar(&logLevel, "log-level", "info", "Log level")
	flag.StringVar(&logFormat, "log-format", "default", "Log format")
	flag.StringVar(&configPath, "config", "config.yaml", "Path to yaml config.")
	flag.StringVar(&metricsPrefix, "prefix", "qb_", "Metrics prefix.")

	setFeatures := feature.Use(useFeatures)
	defer setFeatures()
	flag.Parse()

	config.UpdatePath(configPath)
	metrics.UpdatePrefix(metricsPrefix)
	log.Set(logLevel, logFormat)
}

func main() {
	cfg := config.Get()
	initializeState(cfg)
	client := newHTTPClient(cfg)

	api, err := api.NewQBittorrentAPI(&api.QBittorrentAPIOpts{
		BaseURL: cfg.QBittorrent.BaseURL,
		Credentials: &api.QBittorrentCredentials{
			Username: cfg.QBittorrent.Username,
			Password: cfg.QBittorrent.Password,
		},
		HttpClient: client,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	runScheduledTasks(api, cfg)

	scheduler.Get().Wait()
}

func initializeState(cfg config.Config) {
	if feature.Get(feature.TRANSIENT_STATE) {
		state.SetTransientMode(true)
	} else if cfg.Global.StatePath != "" {
		state.UpdatePath(cfg.Global.StatePath)
	} else {
		log.Debug("No state path configured; state will be transient")
		state.SetTransientMode(true)
	}
}

func newHTTPClient(cfg config.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.QBittorrent.InsecureSkipVerify,
			},
		},
		Timeout: time.Duration(cfg.QBittorrent.Timeout) * time.Second,
	}
}

func runScheduledTasks(api *api.QBittorrentAPI, cfg config.Config) {
	scheduler.Run(func() error {
		http.Handle(cfg.Metrics.UrlPath, promhttp.Handler())
		addr := fmt.Sprintf("http://0.0.0.0:%s%s", cfg.Metrics.Port, cfg.Metrics.UrlPath)
		log.Info("Metrics server is available on port " + addr)
		if err := http.ListenAndServe(":"+cfg.Metrics.Port, nil); err != nil {
			log.Error(err.Error())
		}
		return nil
	}, nil)

	metricsClient := metrics.Get()
	st := state.Get()

	scheduler.Run(func() error {
		torrents, err := api.TorrentsInfo()
		if err != nil {
			return err
		}
		metricsClient.UpdateTorrent(torrents)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: torrentUpdateInterval,
		IsFast:   true,
	})

	scheduler.Run(func() error {
		transfer, err := api.TransferInfo()
		if err != nil {
			return err
		}
		st.UpdateTransferInfo(transfer.DlInfoData, transfer.UpInfoData)
		metricsClient.UpdateTransfer(transfer, st.TransferInfo)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: transferUpdateInterval,
		IsFast:   true,
	})

	scheduler.Run(func() error {
		version, err := api.AppVersion()
		if err != nil {
			return err
		}
		metricsClient.UpdateVersion(version)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: versionCheckInterval,
		IsFast:   true,
	})
}
