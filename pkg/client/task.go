package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/config"
	"github.com/marcosjom/sys-backups-automation/pkg/task"
)

type Task struct {

	// Task configuration file was deleted or failed to load.
	IsOrphaned bool

	// Config
	config config.Task

	// Trigger
	trigger task.Trigger
}

func (t *Task) LoadConfig(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return errors.New("File could not be opened: " + path)
	}
	var payload config.Task
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return errors.New("File could not be unmarshalled: " + path)
	}
	// Set config
	if err := t.trigger.SetConfig(&payload); err != nil {
		return err
	}
	t.config = payload
	//
	return nil
}

func (t *Task) Tick() task.Result {
	r := task.Pending
	now := time.Now()
	// Tick?
	if !t.trigger.IsTickOldEnough(now) {
		return r
	}
	t.trigger.Ticked()
	// Execute?
	if !t.trigger.ShouldRunTask(now) {
		return r
	}
	// Execute
	t.trigger.ExecutionStartedAt(now)
	r = task.Success
	deferred := make([]string, 0, len(t.config.Commands))
	for _, cmdDef := range t.config.Commands {
		cmdStr := strings.TrimSpace(cmdDef.Execute)
		if cmdStr == "" {
			// Empty command is not allowed
			r = task.ErrorUnrecoverable
			break
		}
		cmd := exec.Command(cmdStr)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Command execution failed: %s.\n", err.Error())
			catchStr := strings.TrimSpace(cmdDef.Catch)
			if catchStr != "" {
				cmd2 := exec.Command(catchStr)
				if err := cmd2.Run(); err != nil {
					fmt.Printf("Catch execution failed: %s.\n", err.Error())
				}
			}
			r = task.ErrorRecoverable
			break
		}
		// Append deferred command
		deferredStr := strings.TrimSpace(cmdDef.Deferred)
		if deferredStr != "" {
			deferred = append(deferred, deferredStr)
		}
	}
	//execute deferred commands
	for i := len(deferred) - 1; i >= 0; i-- {
		deferredStr := deferred[i]
		cmd3 := exec.Command(deferredStr)
		if err := cmd3.Run(); err != nil {
			fmt.Printf("Deferred execution failed: %s.\n", err.Error())
		}
	}
	t.trigger.ExecutionEndedAt(time.Now(), r)
	return r
}
