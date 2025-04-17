package types

import "encoding/json"

type Torrent struct {
	AddedOn           int64                      `json:"added_on"`
	AmountLeft        int64                      `json:"amount_left"`
	AutoTMM           bool                       `json:"auto_tmm"`
	Availability      float64                    `json:"availability"`
	Category          string                     `json:"category"`
	Completed         int64                      `json:"completed"`
	CompletionOn      int64                      `json:"completion_on"`
	ContentPath       string                     `json:"content_path"`
	DlLimit           int64                      `json:"dl_limit"`
	Dlspeed           int64                      `json:"dlspeed"`
	Downloaded        int64                      `json:"downloaded"`
	DownloadedSession int64                      `json:"downloaded_session"`
	Eta               int64                      `json:"eta"`
	FLPiecePrio       bool                       `json:"f_l_piece_prio"`
	ForceStart        bool                       `json:"force_start"`
	Hash              string                     `json:"hash"`
	IsPrivate         bool                       `json:"isPrivate"`
	LastActivity      int64                      `json:"last_activity"`
	MagnetURI         string                     `json:"magnet_uri"`
	MaxRatio          float64                    `json:"max_ratio"`
	MaxSeedingTime    int64                      `json:"max_seeding_time"`
	Name              string                     `json:"name"`
	NumComplete       int64                      `json:"num_complete"`
	NumIncomplete     int64                      `json:"num_incomplete"`
	NumLeechs         int64                      `json:"num_leechs"`
	NumSeeds          int64                      `json:"num_seeds"`
	Priority          int64                      `json:"priority"`
	Progress          float64                    `json:"progress"`
	Ratio             float64                    `json:"ratio"`
	RatioLimit        float64                    `json:"ratio_limit"`
	SavePath          string                     `json:"save_path"`
	SeedingTime       int64                      `json:"seeding_time"`
	SeedingTimeLimit  int64                      `json:"seeding_time_limit"`
	SeenComplete      int64                      `json:"seen_complete"`
	SeqDl             bool                       `json:"seq_dl"`
	Size              int64                      `json:"size"`
	State             string                     `json:"state"`
	SuperSeeding      bool                       `json:"super_seeding"`
	Tags              string                     `json:"tags"`
	TimeActive        int64                      `json:"time_active"`
	TotalSize         int64                      `json:"total_size"`
	Tracker           string                     `json:"tracker"`
	UpLimit           int64                      `json:"up_limit"`
	Uploaded          int64                      `json:"uploaded"`
	UploadedSession   int64                      `json:"uploaded_session"`
	Upspeed           int64                      `json:"upspeed"`
	AdditionalFields  map[string]json.RawMessage `json:"-"`
}

type Transfer struct {
	DlInfoSpeed      int64  `json:"dl_info_speed"`
	DlInfoData       int64  `json:"dl_info_data"`
	UpInfoSpeed      int64  `json:"up_info_speed"`
	UpInfoData       int64  `json:"up_info_data"`
	DlRateLimit      int64  `json:"dl_rate_limit"`
	UpRateLimit      int64  `json:"up_rate_limit"`
	DhtNodes         int64  `json:"dht_nodes"`
	ConnectionStatus string `json:"connection_status"`
}
