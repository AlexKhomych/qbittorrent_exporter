package state

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type StateStore struct {
	state State
	path  string
	lock  sync.Mutex
}

var (
	singleInstance *StateStore
)

func NewStateStore(path string) *StateStore {
	if singleInstance != nil {
		slog.Info("StateStore instance is already initialized, skipping...")
		return singleInstance
	}
	return &StateStore{
		path:  path,
		state: readState(path),
	}
}

func GetStateStore() (*StateStore, error) {
	if singleInstance == nil {
		return nil, fmt.Errorf("Please initialize StateStore first")
	}
	return singleInstance, nil
}

type State struct {
	TransferInfo TransferInfoState `json:"transfer_info"`
}

type TransferInfoState struct {
	DlInfoData      int64 `json:"dl_info_data"`
	DlInfoDataTotal int64 `json:"dl_info_data_total"`
	UpInfoData      int64 `json:"up_info_data"`
	UpInfoDataTotal int64 `json:"up_info_data_total"`
}

func (s *StateStore) State() State {
	return s.state
}

func readState(path string) State {
	var state State

	info, err := os.Stat(path)
	if err != nil {
		slog.Error(err.Error() + ", will be created on next write")
		return state
	}
	if info.IsDir() {
		slog.Info("Store path is directory, using empty state")
		return state
	}
	f, err := os.Open(path)
	if err != nil {
		slog.Info(err.Error())
		return state
	}
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		slog.Error(err.Error())
		return state
	}

	return state
}

func (s *StateStore) UpdateTransferInfoState(cDl, cUp int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	state := s.state

	var dDl, dUp int64
	if cDl >= state.TransferInfo.DlInfoData {
		dDl = cDl - state.TransferInfo.DlInfoData
	} else {
		dDl = cDl
		slog.Info(
			"Current DlInfoData is lower than the last recorded one. Posibility of session restart",
			"cur_dl_info_data", cDl,
			"last_dl_info_data", state.TransferInfo.DlInfoData,
		)
	}
	if cUp >= state.TransferInfo.UpInfoData {
		dUp = cUp - state.TransferInfo.UpInfoData
	} else {
		dUp = cUp
		slog.Info(
			"Current UpInfoData is lower than the last recorded one. Posibility of session restart",
			"cur_up_info_data", cDl,
			"last_up_info_data", state.TransferInfo.UpInfoData,
		)
	}

	dlInfoDataTotal := state.TransferInfo.DlInfoDataTotal + dDl
	upInfoDataTotal := state.TransferInfo.UpInfoDataTotal + dUp

	state.TransferInfo.DlInfoData = cDl
	state.TransferInfo.DlInfoDataTotal = dlInfoDataTotal
	state.TransferInfo.UpInfoData = cUp
	state.TransferInfo.UpInfoDataTotal = upInfoDataTotal

	s.state = state

	f, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("Failed to open state store file")
	}
	if err := json.NewEncoder(f).Encode(&s.state); err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("Failed to encode json to a state store file")
	}

	return nil
}
