package main

import (
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

const version = "1.0.1"

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

		useFeatures = map[feature.FeatureFlag]bool{
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
	config := config.Get()
	state.UpdatePath(config.Global.StatePath)

	api, err := api.NewQBittorrentAPI(&api.QBittorrentAPIOpts{
		BaseURL: config.QBittorrent.BaseURL,
		Credentials: &api.QBittorrentCredentials{
			Username: config.QBittorrent.Username,
			Password: config.QBittorrent.Password,
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	runInParallel(api, config)
	scheduler.Get().Wait()
}

func runInParallel(api *api.QBittorrentAPI, config config.Config) {
	scheduler.Run(func() error {
		http.Handle(config.Metrics.UrlPath, promhttp.Handler())
		log.Info("Metrics server is available on port http://0.0.0.0:" + config.Metrics.Port + config.Metrics.UrlPath)
		if err := http.ListenAndServe(":"+config.Metrics.Port, nil); err != nil {
			log.Error(err.Error())
		}
		return nil
	}, nil)

	metrics := metrics.Get()
	state := state.Get()

	scheduler.Run(func() error {
		torrents, err := api.TorrentsInfo()
		if err != nil {
			return err
		}
		metrics.UpdateTorrent(torrents)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: 30 * time.Second,
		IsFast:   true,
	})

	scheduler.Run(func() error {
		transfer, err := api.TransferInfo()
		if err != nil {
			return err
		}
		state.UpdateTransferInfo(transfer.DlInfoData, transfer.UpInfoData)
		metrics.UpdateTransfer(transfer, state.TransferInfo)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: 30 * time.Second,
		IsFast:   true,
	})

	scheduler.Run(func() error {
		verison, err := api.AppVersion()
		if err != nil {
			return err
		}
		metrics.UpdateVersion(verison)
		return nil
	}, &scheduler.PeriodicTaskOpts{
		Interval: 10 * time.Minute,
		IsFast:   true,
	})
}
