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

	// Defines the minimun age between attemps and executions.
	Age struct {
		// Minimun time between evaluations of the trigger.
		Tick Age //years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		// Minimun time between executions.
		Min Age //years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		// If defined, an execution is attempted even outside of the allowed range of instans.
		// If not defined, a succesful execution is allowed to be skipped if wa not possible inside the allowed's range.
		Max Age //years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
	}
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
