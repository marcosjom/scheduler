package client

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Config struct {

	// Configurations (read only)

	Configs struct {

		// Folderpath with tasks ".json" files ("".state.json")
		Path string

		// Seconds between folder read; for new/updated tasks detection.
		SecsBetweenSync int
	}

	State struct {

		// Seconds between dumping of persistent-state to file.
		SecsBetweenSave int
	}
}

// Values validation.
func (c *Config) HasError() error {
	// Validate path
	if strings.TrimSpace(c.Configs.Path) == "" {
		return errors.New("configs.path is required.")
	}
	// Validate secsBetweenSync
	if c.Configs.SecsBetweenSync <= 0 {
		return errors.New("configs.secsBetweenSync is required and must be positive.")
	}
	//
	return nil
}

func (c *Config) Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	var payload Config
	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	//
	*c = payload
	//
	return nil
}
