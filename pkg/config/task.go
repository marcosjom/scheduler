package config

// Defines a task to be executed.
type Task struct {

	// The timing for triggering this task.
	Timing Timing

	// The commands for the task.
	Commands []Command
}

// Lexical validation.
func (t Task) HasError() error {
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
