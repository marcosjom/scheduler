package client

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcosjom/sys-backups-automation/pkg/task"
)

type Client struct {

	// Configuration
	config Config

	// Tasks
	tasks struct {

		// Array/slice of tasks
		arr map[string]*Task

		// Seconds since tasks configuration files were synced.
		secsSinceSync uint32
	}

	// States
	state struct {
		filepath        string
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

func (c *Client) Prepare(configFilepath string, stateFilepath string) error {
	// config
	config := Config{}
	if err := config.Load(configFilepath); err != nil {
		return err
	}
	if err := config.HasError(); err != nil {
		return err
	}
	// persistent-state
	state := State{}
	if err := state.Load(stateFilepath); err != nil {
		return err
	}
	// load tasks
	if err := c.reloadTasksWithState(config.Configs.Path, &state); err != nil {
		return err
	}
	// set
	c.config = config
	c.state.filepath = stateFilepath
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
				fmt.Printf("Load failed for: '%s': '%s'.\n", filepath, err.Error())
			} else {
				fmt.Printf("Loaded: '%s'.\n", filepath)
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
				// Determine if state changed
				if !fndSavedState || savedState.VersionId != task.trigger.History.VersionId {
					c.state.hasChanged = true
				}
			}
		} else {
			fmt.Printf("Ignoring: '%s'.\n", fileName)
		}
	}
	// Remove all orphaned tasks
	for k, t := range c.tasks.arr {
		if t.IsOrphaned {
			delete(c.tasks.arr, k)
		}
	}
	//
	fmt.Printf("%d -> %d tasks remain after syncing.\n", tasksCountBefore, len(c.tasks.arr))
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
	if err := c.state.payload.Save(c.state.filepath); err != nil {
		return err
	}
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
		result := t.Tick()
		if result != task.Pending {
			c.state.hasChanged = true
		}
	}
	//
	return nil
}
