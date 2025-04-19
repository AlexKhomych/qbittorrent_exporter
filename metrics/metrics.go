package metrics

import (
	"fmt"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/state"
	"qbittorrent_exporter/types"
	"reflect"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	lock                  = &sync.Mutex{}
	metricsPrefix  string = "qb_"
	singleInstance *Metrics
)

type Metrics struct {
	torrent  *torrentMetrics
	transfer *transferMetrics
	version  *versionMetrics
}

type torrentMetrics struct {
	Name       *prometheus.GaugeVec
	State      *prometheus.GaugeVec
	Progress   *prometheus.GaugeVec
	DlSpeed    *prometheus.GaugeVec
	UpSpeed    *prometheus.GaugeVec
	Downloaded *prometheus.GaugeVec
	AmountLeft *prometheus.GaugeVec
	Ratio      *prometheus.GaugeVec
	Eta        *prometheus.GaugeVec
	NumSeeds   *prometheus.GaugeVec
	NumLeechs  *prometheus.GaugeVec
}

type transferMetrics struct {
	Status          *prometheus.GaugeVec
	DlInfoSpeed     *prometheus.GaugeVec
	DlInfoData      *prometheus.GaugeVec
	UpInfoSpeed     *prometheus.GaugeVec
	UpInfoData      *prometheus.GaugeVec
	DlRateLimit     *prometheus.GaugeVec
	UpRateLimit     *prometheus.GaugeVec
	DhtNodes        *prometheus.GaugeVec
	DlInfoDataTotal *prometheus.GaugeVec
	UpInfoDataTotal *prometheus.GaugeVec
}

type versionMetrics struct {
	Version *prometheus.GaugeVec
}

func UpdatePrefix(prefix string) {
	metricsPrefix = prefix
}

func Get() *Metrics {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		singleInstance = &Metrics{}
		singleInstance.initialize()
	}
	return singleInstance
}

func (m *Metrics) initialize() {
	m.torrent = &torrentMetrics{
		Name: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_name",
			Help: "Name of the torrent",
		}, []string{"name"}),

		State: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_state",
			Help: "State of the torrent",
		}, []string{"name", "state"}),

		Progress: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_progress",
			Help: "Progress of the torrent",
		}, []string{"name"}),

		DlSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_dlspeed",
			Help: "Download speed of the torrent",
		}, []string{"name"}),

		UpSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_upspeed",
			Help: "Upload speed of the torrent",
		}, []string{"name"}),

		Downloaded: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_downloaded",
			Help: "Amount of data downloaded",
		}, []string{"name"}),

		AmountLeft: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_amount_left",
			Help: "Amount of data left to download",
		}, []string{"name"}),

		Ratio: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_ratio",
			Help: "Torrent share ratio",
		}, []string{"name"}),

		Eta: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_eta",
			Help: "Estimated time to completion",
		}, []string{"name"}),

		NumSeeds: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_num_seeds",
			Help: "Number of seeds connected to",
		}, []string{"name"}),

		NumLeechs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "torrent_num_leechs",
			Help: "Number of leechers connected to",
		}, []string{"name"}),
	}

	m.transfer = &transferMetrics{
		Status: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_connection_status",
			Help: "Connection status",
		}, []string{}),

		DlInfoSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_dl_info_speed",
			Help: "Global download rate (bytes/s)",
		}, []string{}),

		DlInfoData: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_dl_info_data",
			Help: "Data downloaded this session (bytes)",
		}, []string{}),

		UpInfoSpeed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_up_info_speed",
			Help: "Global upload rate (bytes/s)",
		}, []string{}),

		UpInfoData: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_up_info_data",
			Help: "Data uploaded this session (bytes)",
		}, []string{}),

		DlRateLimit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_dl_rate_limit",
			Help: "Download rate limit (bytes/s)",
		}, []string{}),

		UpRateLimit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_up_rate_limit",
			Help: "Upload rate limit (bytes/s)",
		}, []string{}),

		DhtNodes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_dht_nodes",
			Help: "DHT nodes connected to",
		}, []string{}),

		DlInfoDataTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_dl_info_data_total",
			Help: "Data downloaded total (bytes)",
		}, []string{}),

		UpInfoDataTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "transfer_up_info_data_total",
			Help: "Data downloaded total (bytes)",
		}, []string{}),
	}

	m.version = &versionMetrics{
		Version: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsPrefix + "app_version",
			Help: "Application version",
		}, []string{"version"}),
	}

	registerMetrics(m.torrent)
	registerMetrics(m.transfer)
	registerMetrics(m.version)
}

func (m *Metrics) UpdateTorrent(torrents []types.Torrent) {
	tm := m.torrent
	for _, torrent := range torrents {
		tm.Name.WithLabelValues(torrent.Name).Set(1)
		tm.State.WithLabelValues(torrent.Name, torrent.State).Set(1)
		tm.Progress.WithLabelValues(torrent.Name).Set(torrent.Progress)
		tm.DlSpeed.WithLabelValues(torrent.Name).Set(float64(torrent.Dlspeed))
		tm.UpSpeed.WithLabelValues(torrent.Name).Set(float64(torrent.Upspeed))
		tm.Downloaded.WithLabelValues(torrent.Name).Set(float64(torrent.Downloaded))
		tm.AmountLeft.WithLabelValues(torrent.Name).Set(float64(torrent.AmountLeft))
		tm.Ratio.WithLabelValues(torrent.Name).Set(float64(torrent.Ratio))
		tm.Eta.WithLabelValues(torrent.Name).Set(float64(torrent.Eta))
		tm.NumSeeds.WithLabelValues(torrent.Name).Set(float64(torrent.NumSeeds))
		tm.NumLeechs.WithLabelValues(torrent.Name).Set(float64(torrent.NumLeechs))
	}
}

func (m *Metrics) UpdateTransfer(transfer types.Transfer, state state.TransferInfoState) {
	tm := m.transfer
	var status float64 = 0
	if transfer.ConnectionStatus == "connected" {
		status = 1
	}
	tm.Status.WithLabelValues().Set(status)
	tm.DlInfoSpeed.WithLabelValues().Set(float64(transfer.DlInfoSpeed))
	tm.DlInfoData.WithLabelValues().Set(float64(transfer.DlInfoData))
	tm.UpInfoSpeed.WithLabelValues().Set(float64(transfer.UpInfoSpeed))
	tm.UpInfoData.WithLabelValues().Set(float64(transfer.UpInfoData))
	tm.DlRateLimit.WithLabelValues().Set(float64(transfer.DlRateLimit))
	tm.UpRateLimit.WithLabelValues().Set(float64(transfer.UpRateLimit))
	tm.DhtNodes.WithLabelValues().Set(float64(transfer.DhtNodes))

	log.Debug(fmt.Sprintf("UpdateTransfer() call; DlInfoDataTotal: %d", state.DlInfoDataTotal))
	log.Debug(fmt.Sprintf("UpdateTransfer() call; UpInfoDataTotal: %d", state.UpInfoDataTotal))
	tm.DlInfoDataTotal.WithLabelValues().Set(float64(state.DlInfoDataTotal))
	tm.UpInfoDataTotal.WithLabelValues().Set(float64(state.UpInfoDataTotal))
}

func (m *Metrics) UpdateVersion(version string) {
	vm := m.version
	vm.Version.WithLabelValues(version).Set(1)
}

// registerMetrics accepts MetricsStruct
// which contains multiple metrics fields
func registerMetrics(metrics any) {
	val := reflect.ValueOf(metrics)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		log.Error("Expected a struct, got " + val.Kind().String())
	}

	typ := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			log.Debug("Skipping unexported field", "name", fieldType.Name)
			continue
		}

		if field.Kind() == reflect.Ptr && !field.IsNil() {
			if metric, ok := field.Interface().(prometheus.Collector); ok {
				prometheus.MustRegister(metric)
				log.Debug("Metrics was registered", "name", fieldType.Name)
			}
		}
	}
}
