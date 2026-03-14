package client

import (
	"bytes"
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
	cmd := exec.Command("bash", "-c", s)
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
		outb := bytes.Buffer{}
		errb := bytes.Buffer{}
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			if outStr := outb.String(); len(outStr) > 0 {
				log.Printf("Command stdout: %s.\n", outStr)
			}
			if errStr := errb.String(); len(errStr) > 0 {
				log.Printf("Command stderr: %s.\n", errStr)
			}
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
				outb2 := bytes.Buffer{}
				errb2 := bytes.Buffer{}
				cmd2.Stdout = &outb2
				cmd2.Stderr = &errb2
				if err := cmd2.Run(); err != nil {
					if outStr2 := outb2.String(); len(outStr2) > 0 {
						log.Printf("Catch stdout: %s.\n", outStr2)
					}
					if errStr2 := errb2.String(); len(errStr2) > 0 {
						log.Printf("Catch stderr: %s.\n", errStr2)
					}
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
		outb3 := bytes.Buffer{}
		errb3 := bytes.Buffer{}
		cmd3.Stdout = &outb3
		cmd3.Stderr = &errb3
		if err := cmd3.Run(); err != nil {
			if outStr3 := outb3.String(); len(outStr3) > 0 {
				log.Printf("Deferred stdout: %s.\n", outStr3)
			}
			if errStr3 := errb3.String(); len(errStr3) > 0 {
				log.Printf("Deferred stderr: %s.\n", errStr3)
			}
			log.Printf("Deferred execution failed: %s.\n", err.Error())
		}
	}
	t.trigger.ExecutionEndedAt(time.Now(), r)
	t.ExecSeq++
	log.Printf("Task executed: '%s'.\n", t.path)
	return r
}
