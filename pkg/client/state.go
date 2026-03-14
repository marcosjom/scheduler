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

func NewState() State {
	return State{Tasks: make(map[string]task.History)}
}

func (s *State) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Does not exists
			return false, nil
		}
		// Error
		return false, err
	}
	// Exists
	return true, nil
}

func (s *State) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	payload := NewState()
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	// Validate emty Tasks
	if payload.Tasks == nil {
		payload.Tasks = make(map[string]task.History)
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
