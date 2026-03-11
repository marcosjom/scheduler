package config

import (
	"encoding/json"
	"hash/crc32"
)

// Defines a task to be executed.
type Task struct {

	// Task configuration version.
	// Only necesary of you want to trigger a configuration change event, without modifying the configuration.
	Version int

	// The timing for triggering this task.
	Timing Timing

	// The commands for the task.
	Commands []Command
}

// Lexical validation.
func (t *Task) HasError() error {
	//Timing
	if err := t.Timing.HasError(); err != nil {
		return err
	}
	//Cmds
	for _, cmd := range t.Commands {
		if err := cmd.HasError(); err != nil {
			return err
		}
	}
	return nil
}

// Task's config CRC32; to detect changes on the configuration.

func (t *Task) Crc32() uint32 {
	b, err := json.Marshal(t)
	if err != nil {
		return 0
	}
	table := crc32.MakeTable(0)
	return crc32.Checksum(b, table)
}
