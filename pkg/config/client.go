package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Client struct {

	// Configurations (read only)

	Configs struct {

		// Folderpath with tasks ".json" files to read.
		Path string

		// Seconds between folder read; for new/updated tasks detection.
		SecsBetweenSync int
	}

	State struct {

		// Folderpath with persistent-state file to write.
		Path string

		// Seconds between dumping of persistent-state to file.
		SecsBetweenSave int
	}
}

// Values validation.
func (c *Client) HasError() error {
	// Validate path
	if strings.TrimSpace(c.Configs.Path) == "" {
		return errors.New("configs.path is required.")
	}
	// Validate secsBetweenSync
	if c.Configs.SecsBetweenSync <= 0 {
		return errors.New("configs.secsBetweenSync is required and must be positive.")
	}
	// Validate path
	if strings.TrimSpace(c.State.Path) == "" {
		return errors.New("state.path is required.")
	}
	//
	return nil
}

func (c *Client) Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	var payload Client
	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	//
	*c = payload
	//
	return nil
}
