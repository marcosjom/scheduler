package client

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcosjom/scheduler/pkg/config"
)

type Client struct {

	// Configuration
	config config.Client

	// Tasks
	tasks struct {

		// Array/slice of tasks
		arr map[string]*Task

		// Seconds since tasks configuration files were synced.
		secsSinceSync uint32
	}

	// States
	state struct {
		payload         State
		hasChanged      bool
		secsWithoutSave int
	}
}

// Prepares the client

func New() *Client {
	r := Client{}
	r.tasks.arr = make(map[string]*Task)
	return &r
}

func (c *Client) Prepare(config config.Client) error {
	// config
	if err := config.HasError(); err != nil {
		return err
	}
	// persistent-state
	state := NewState()
	stateExists, stateErr := state.FileExists(config.State.Path)
	if stateErr != nil {
		return stateErr
	}
	// Check for readable-state
	if stateExists {
		if err := state.Load(config.State.Path); err != nil {
			return err
		}
	}
	// Load tasks
	if err := c.reloadTasksWithState(config.Configs.Path, &state); err != nil {
		return err
	}
	// Check for writable-state
	if err := state.Save(config.State.Path); err != nil {
		return err
	}
	// set
	c.config = config
	c.state.payload = state
	//
	return nil
}

// Reloads the tasks from json files inside the specified folder.

func (c *Client) ReloadTasks(folderPathP string) error {
	return c.reloadTasksWithState(folderPathP, &c.state.payload)
}

// Reloads the tasks from json files inside the specified folder.
func (c *Client) reloadTasksWithState(folderPathP string, state *State) error {
	tasksCountBefore := len(c.tasks.arr)
	// Reset sync counter
	if folderPathP == c.config.Configs.Path {
		c.tasks.secsSinceSync = 0
	}
	// Determine folderpath
	folderPath := strings.TrimSpace(folderPathP)
	for strings.HasSuffix(folderPath, "/") || strings.HasSuffix(folderPath, "\\") {
		folderPath = folderPath[:len(folderPath)-1]
	}
	if folderPath == "" {
		return errors.New("Empty folderPath.")
	}
	// Read folder
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return errors.New("Could not read directory.")
	}
	// Flag all tasks as Orphaned
	for _, t := range c.tasks.arr {
		t.IsOrphaned = true
	}
	// Load files
	for _, e := range entries {
		fileName := e.Name()
		fileExt := filepath.Ext(fileName)
		if strings.EqualFold(fileExt, ".json") {
			filepath := folderPath + string(filepath.Separator) + fileName
			// Find task
			task, fnd := c.tasks.arr[filepath]
			if !fnd || task == nil {
				task = &Task{IsOrphaned: true}
				fnd = false
			}
			// Load file
			if err := task.LoadConfig(filepath); err != nil {
				log.Printf("Load failed for: '%s': '%s'.\n", filepath, err.Error())
			} else {
				task.IsOrphaned = false
				savedState, fndSavedState := state.Tasks[filepath]
				// Add to array
				if !fnd {
					c.tasks.arr[filepath] = task
					if fndSavedState {
						// Apply saved persistent-state
						task.trigger.History = savedState
					} else {
						// Flag as changed
						task.trigger.History.VersionId++
					}
				}
			}
		} else {
			log.Printf("Ignoring: '%s'.\n", fileName)
		}
	}
	// Remove all orphaned tasks
	for k, t := range c.tasks.arr {
		if t.IsOrphaned {
			delete(c.tasks.arr, k)
		}
	}
	//
	log.Printf("%d -> %d tasks remain after syncing.\n", tasksCountBefore, len(c.tasks.arr))
	//
	return nil
}

// Save state
func (c *Client) SaveState() error {
	c.state.secsWithoutSave = 0
	// sync tasks
	for k, t := range c.tasks.arr {
		c.state.payload.Tasks[k] = t.trigger.History
	}
	// dump to file
	if err := c.state.payload.Save(c.config.State.Path); err != nil {
		return err
	}
	log.Printf("Client-state saved.\n")
	//
	return nil
}

// Client's logic heartbeat

func (c *Client) TickOneSecond() error {
	// Save changes
	if c.state.hasChanged {
		c.state.secsWithoutSave++
		if c.state.secsWithoutSave >= c.config.State.SecsBetweenSave {
			c.state.secsWithoutSave = 0
			c.state.hasChanged = false
			c.SaveState()
		}
	}
	// Synchronize tasks from folder
	c.tasks.secsSinceSync++
	secsBetweenSync := uint32(c.config.Configs.SecsBetweenSync)
	if secsBetweenSync > 0 && c.tasks.secsSinceSync >= secsBetweenSync {
		// Reset sync counter
		c.tasks.secsSinceSync = 0
		c.ReloadTasks(c.config.Configs.Path)
	}
	// Evaluate tasks
	for _, t := range c.tasks.arr {
		t.Tick()
		// Detect unsaved state changes
		if t.ExecSeqSaved != t.ExecSeq {
			t.ExecSeqSaved = t.ExecSeq
			c.state.hasChanged = true
		}
	}
	//
	return nil
}
