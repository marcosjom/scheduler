package config

type Timing struct {

	// Defines the range on which the action should triggered.
	Range struct {
		// Examples are:
		// For yearly tasks: min: "december 31", max "december 31"
		// For monthly tasks: min: "1", max "1"
		// For weekly tasks: min: "monday 00h 00m", max "monday 00h 59m"
		// For daily tasks: min: "00h 00m", max "00h 59m"
		Min Instant //month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
		Max Instant //month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
	}

	// Defines the minimun and maximun wait time since last
	// execution of the task associated to the trigger.
	Age struct {
		Min Age //years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		Max Age //years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
	}

	// If true, the task is not allowed to run out of its specified Range, event if last execution is over Ange.Max old.
	// If false, the task ignores the specified Range in attempt to avoid executions older than Ange.Max.
	IsSkipAllowed bool

	// If true, multiple tasks can run simultaniously; usually when the task takes longer than its Range.
	// If false, only one instance of this task can be active.
	IsMultipleInstancesAllowed bool
}

// Lexical validation.
func (t *Timing) HasError() error {
	// Range
	if err := t.Range.Min.HasError(); err != nil {
		return err
	}
	if err := t.Range.Max.HasError(); err != nil {
		return err
	}
	// Age
	if err := t.Age.Min.HasError(); err != nil {
		return err
	}
	if err := t.Age.Max.HasError(); err != nil {
		return err
	}
	//
	return nil
}
