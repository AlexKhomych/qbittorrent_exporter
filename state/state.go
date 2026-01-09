package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"qbittorrent_exporter/lib/log"
	"qbittorrent_exporter/lib/scheduler"
	"qbittorrent_exporter/validator"
	"sync"
	"time"
)

var (
	statePath      = "state.json"
	lock           sync.Mutex
	singleInstance *State
	transientMode  = false
)

type State struct {
	TransferInfo TransferInfoState `json:"transfer_info"`
}

type TransferInfoState struct {
	DlInfoData      int64 `json:"dl_info_data"`
	DlInfoDataTotal int64 `json:"dl_info_data_total"`
	UpInfoData      int64 `json:"up_info_data"`
	UpInfoDataTotal int64 `json:"up_info_data_total"`
}

func init() {
	scheduler.Run(func() error {
		lock.Lock()
		defer lock.Unlock()
		if transientMode || singleInstance == nil {
			return nil
		}
		return singleInstance.write()
	}, &scheduler.PeriodicTaskOpts{
		Interval: 30 * time.Second,
		IsFast:   false,
	})
}

func UpdatePath(path string) {
	lock.Lock()
	defer lock.Unlock()
	statePath = path
}

func Get() *State {
	lock.Lock()
	defer lock.Unlock()
	if singleInstance == nil {
		if transientMode {
			singleInstance = &State{}
		} else {
			singleInstance = readState(statePath)
		}
	}
	return singleInstance
}

func SetTransientMode(isTransient bool) {
	lock.Lock()
	defer lock.Unlock()
	transientMode = isTransient
	if isTransient {
		log.Info("State set to transient mode; no state will be persisted")
		singleInstance = &State{}
	} else {
		singleInstance = readState(statePath)
	}
}

func readState(path string) *State {
	if transientMode {
		return &State{}
	}

	var state State
	if err := validator.ValidatePath(path, false); err != nil {
		log.Warn(err.Error())
		log.Info("State file will be created on the next write")
		return &State{}
	}

	f, err := os.Open(path)
	if err != nil {
		log.Error(err.Error() + ". Using transient state")
		return &State{}
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&state); err != nil {
		log.Error(err.Error())
		return &State{}
	}

	return &state
}

func (s *State) UpdateTransferInfo(dl, up int64) {
	lock.Lock()
	defer lock.Unlock()
	s.TransferInfo.calculateDelta(dl, up)
}

func (s *State) write() error {
	if transientMode {
		return nil
	}

	absPath, err := filepath.Abs(statePath)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug("Writing state into a file: " + absPath)
	log.Debug(fmt.Sprintf("State: %+v", s))

	f, err := os.OpenFile(statePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to open state store file")
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&s); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("Failed to encode json to a state store file")
	}

	return nil
}

func (t *TransferInfoState) calculateDelta(dl, up int64) {
	delta := func(current, previous int64) int64 {
		if current >= previous {
			return current - previous
		}
		log.Info("Current value is lower than the previously recorded one. Possibility of session restart")
		return current
	}

	t.DlInfoDataTotal += delta(dl, t.DlInfoData)
	t.UpInfoDataTotal += delta(up, t.UpInfoData)
	t.DlInfoData = dl
	t.UpInfoData = up

	log.Debug(fmt.Sprintf("TransferInfoState update: %+v", t))
}
