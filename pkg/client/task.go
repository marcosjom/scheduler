package client

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/config"
	"github.com/marcosjom/sys-backups-automation/pkg/task"
)

type Task struct {

	// Task configuration file was deleted or failed to load.
	IsOrphaned bool

	// Current execution sequence
	ExecSeq uint32

	// Current execution saved
	ExecSeqSaved uint32

	// Config
	config config.Task
	path   string

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
	t.path = path
	//
	return nil
}

func argsFromString(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' ' // space
	fields, err := r.Read()
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func cmdFromString(s string) (*exec.Cmd, error) {
	cmdParts, cmdErr := argsFromString(s)
	if cmdErr != nil {
		return nil, errors.New("Cmd-args parsing failed.")
	}
	//
	os := runtime.GOOS
	if os == "windows" {
		//add: "cmd", "/C"
		cmd := exec.Command("cmd")
		cmd.Args = append([]string{"/C"}, cmdParts...)
		return cmd, nil
	}
	cmd := exec.Command(cmdParts[0])
	if len(cmdParts) > 1 {
		cmd.Args = cmdParts[1:]
	}
	return cmd, nil
}

func (t *Task) Tick() task.Result {
	r := task.Pending
	now := time.Now()
	// Should Evaluate?
	if !t.trigger.IsTickOldEnough(now) {
		return r
	}
	t.trigger.Ticked()
	// Should Execute?
	if !t.trigger.ShouldRunTask(now) {
		return r
	}
	// Execute
	t.trigger.ExecutionStartedAt(now)
	t.ExecSeq++
	r = task.Success
	deferred := make([]string, 0, len(t.config.Commands))
	for _, cmdDef := range t.config.Commands {
		cmdStr := strings.TrimSpace(cmdDef.Execute)
		if cmdStr == "" {
			// Empty command is not allowed
			log.Printf("Task unrecoverable-error, empty command.\n")
			r = task.ErrorUnrecoverable
			break
		}
		cmd, cmdErr := cmdFromString(cmdStr)
		if cmdErr != nil {
			// Args are invalid
			log.Printf("Task unrecoverable-error, cmd-args parsing failed: '%s'.\n", cmdStr)
			r = task.ErrorUnrecoverable
			break
		}
		if err := cmd.Run(); err != nil {
			log.Printf("Command execution failed: %s.\n", err.Error())
			catchStr := strings.TrimSpace(cmdDef.Catch)
			if catchStr != "" {
				cmd2, cmd2Err := cmdFromString(catchStr)
				if cmd2Err != nil {
					// Args are invalid
					log.Printf("Task unrecoverable-error, catch-args parsing failed: '%s'.\n", catchStr)
					r = task.ErrorUnrecoverable
					break
				}
				if err := cmd2.Run(); err != nil {
					log.Printf("Catch execution failed: %s.\n", err.Error())
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
		cmd3, cmd3Err := cmdFromString(deferredStr)
		if cmd3Err != nil {
			// Args are invalid
			log.Printf("Task unrecoverable-error, deferred-args parsing failed: '%s'.\n", deferredStr)
			r = task.ErrorUnrecoverable
			break
		}
		if err := cmd3.Run(); err != nil {
			log.Printf("Deferred execution failed: %s.\n", err.Error())
		}
	}
	t.trigger.ExecutionEndedAt(time.Now(), r)
	t.ExecSeq++
	log.Printf("Task executed: '%s'.\n", t.path)
	return r
}
