package config

import (
	"errors"
	"strings"
)

// Defined a command to be executed.
type Command struct {

	// Command to be executed.
	Execute string

	// Command to be executed if 'Execute' fails (optional).
	Catch string

	// Command to be executed if 'Execute' succeds; as clenup reverse order (like deleting a temporary file).
	Deferred string
}

// Lexical validation.
func (t *Command) HasError() error {
	// Execute
	if strings.TrimSpace(t.Execute) == "" {
		return errors.New("Command.Execute cannot be empty.")
	}
	//
	return nil
}
