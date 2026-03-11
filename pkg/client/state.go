package client

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/marcosjom/sys-backups-automation/pkg/task"
)

type State struct {
	Tasks map[string]task.History
}

func (s *State) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	var payload State
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	// Set config
	*s = payload
	//
	return nil
}

func (s *State) Save(path string) error {
	bytes, _ := json.Marshal(s)
	err := os.WriteFile(path, bytes, 0644)
	if err != nil {
		return errors.New("File could not be saved: " + path)
	}
	//
	return nil
}
