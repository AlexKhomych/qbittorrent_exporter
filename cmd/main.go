package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"qbittorrent_exporter/config"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/lib/qbittorrent/api"
	"qbittorrent_exporter/lib/scheduler"
	"qbittorrent_exporter/metrics"
	"qbittorrent_exporter/state"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version = "v0.1.0-alpha"

func main() {
	config := config.Get().WithEnvPriority()
	var wg sync.WaitGroup

	taskMetricServer := &SpawnTaskOpts{
		wg: &wg,
		task: func() error {
			http.Handle(config.Metrics.UrlPath, promhttp.Handler())
			if err := http.ListenAndServe(":"+config.Metrics.Port, nil); err != nil {
				log.Error(err.Error())
			} else {
				log.Info("Metrics server is available on port http://0.0.0.0:" + config.Metrics.Port + config.Metrics.UrlPath)
			}
			return nil
		},
	}
	go SpawnTask(taskMetricServer)

	api, err := api.NewQBittorrentAPI(&api.QBittorrentAPIOpts{
		BaseURL: config.QBittorrent.BaseURL,
		Credentials: &api.QBittorrentCredentials{
			Username: config.QBittorrent.Username,
			Password: config.QBittorrent.Password,
		},
	})
	check(err)

	metrics := metrics.Get()
	ss := state.NewStateStore(stateStorePath)

	taskUpdateTorrents := &SpawnTaskOpts{
		wg: &wg,
		task: func() error {
			torrents, err := api.TorrentsInfo()
			if err != nil {
				return err
			}
			metrics.UpdateTorrent(torrents)
			return nil
		},
		po: &scheduler.PeriodicTaskOpts{
			Interval: 30 * time.Second,
			IsFast:   true,
		},
	}
	go SpawnTask(taskUpdateTorrents)

	taskUpdateTransfer := &SpawnTaskOpts{
		wg: &wg,
		task: func() error {
			transfer, err := api.TransferInfo()
			if err != nil {
				return err
			}

			if err := ss.UpdateTransferInfoState(transfer.DlInfoData, transfer.UpInfoData); err != nil {
				return err
			}
			state := ss.State().TransferInfo
			metrics.UpdateTransfer(transfer, state)
			return nil
		},
		po: &scheduler.PeriodicTaskOpts{
			Interval: 30 * time.Second,
			IsFast:   true,
		},
	}
	go SpawnTask(taskUpdateTransfer)

	taskUpdateVersion := &SpawnTaskOpts{
		wg: &wg,
		task: func() error {
			verison, err := api.AppVersion()
			if err != nil {
				return err
			}
			metrics.UpdateVersion(verison)
			return nil
		},
		po: &scheduler.PeriodicTaskOpts{
			Interval: 10 * time.Minute,
			IsFast:   true,
		},
	}
	go SpawnTask(taskUpdateVersion)

	wg.Wait()
}

type SpawnTaskOpts struct {
	wg   *sync.WaitGroup
	po   *scheduler.PeriodicTaskOpts
	task func() error
}

func SpawnTask(o *SpawnTaskOpts) {
	o.wg.Add(1)
	defer o.wg.Done()
	if o.po != nil {
		scheduler.RunPeriodicTask(o.task, o.po)
	} else if err := o.task(); err != nil {
		log.Error(err.Error())
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	logFormat string
	logLevel  string

	stateStorePath string
	configPath     string
)

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintf(w, "\nVersion: %v\n\nUsage: %v [ Options... ]\n\nAvailable Options:\n",
			version, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.StringVar(&logFormat, "log-format", "default", "Log format")
	flag.StringVar(&logLevel, "log-level", "info", "Log level")
	flag.StringVar(&stateStorePath, "state-store-path", "state.json", "Path for state storage")
	flag.StringVar(&configPath, "config-path", "config.yaml", "Path to yaml config")

	flag.Parse()

	config.UpdatePath(configPath)
	log.InitializeLog(logFormat, logLevel)
}
