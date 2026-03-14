package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Scheduler struct {

	// Client
	Client Client
}

// Values validation.
func (c *Scheduler) HasError() error {
	// Client
	if err := c.Client.HasError(); err != nil {
		return err
	}
	//
	return nil
}

func (c *Scheduler) Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	var payload Scheduler
	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	//
	*c = payload
	//
	return nil
}
