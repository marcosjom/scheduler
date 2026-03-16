package task

import (
	"testing"
	"time"

	"github.com/marcosjom/scheduler/pkg/config"
)

// Testing calculating incremental age, one second increments during one year.
func Test_Trigger_IsTickOldEnough(t *testing.T) {
	start := time.Now()
	hoursMaxToEval := 128
	// Prepare trigger
	config := config.Task{}
	{
		config.Timing.Age.Tick = "1d"
		config.Timing.Age.Min = "23h"
		config.Timing.Age.Max = "25h"
		//
		config.Timing.Range.Min = "January"
		config.Timing.Range.Max = "January"
	}
	trigger := Trigger{}
	if err := trigger.SetConfig(&config); err != nil {
		t.Errorf("SetConfig failed: %s.", err.Error())
		return
	}
	// Test
	execsCount := 0
	for iHour := 0; iHour <= hoursMaxToEval; iHour++ {
		end := start.Add(time.Hour * time.Duration(iHour))
		if trigger.IsTickOldEnough(end) {
			trigger.History.LastTick.Time = end
			execsCount++
		}
	}
	expectedExecCount := 1 + (hoursMaxToEval / 24)
	if execsCount != expectedExecCount {
		t.Errorf("Expected %d ticks (%d done) for '%s'-tick in %d hours.", expectedExecCount, execsCount, config.Timing.Age.Tick, hoursMaxToEval)
		return
	}
}

// Testing calculating task monthly execution during a year (tick per day).
func Test_Trigger_ShouldRunTask_Daily(t *testing.T) {
	start := time.Now()
	monthsToEval := 12
	daysMaxToEval := (monthsToEval * 31)
	hoursMaxToEval := (daysMaxToEval * 24)
	// Prepare trigger
	config := config.Task{}
	{
		config.Timing.Age.Tick = "6h"
		config.Timing.Age.Min = "23h"
		//config.Timing.Age.Max = "25h" //no too-old-value
		//
		config.Timing.Range.Min = "1" //1 of each month
		config.Timing.Range.Max = "3" //3 of each month
	}
	trigger := Trigger{}
	if err := trigger.SetConfig(&config); err != nil {
		t.Errorf("SetConfig failed: %s.", err.Error())
		return
	}
	// Test
	ticksCount := 0
	execsCount := 0
	for iHour := 0; iHour <= hoursMaxToEval; iHour++ {
		end := start.Add(time.Hour * time.Duration(iHour))
		//Tick
		if trigger.IsTickOldEnough(end) {
			trigger.History.LastTick.Time = end
			ticksCount++
			//Run
			if trigger.ShouldRunTask(end) {
				//log.Printf("Executing: %s.\n", end.String()[:16])
				trigger.History.LastSuccess.Time = end
				execsCount++
			}
		}
	}
	expectedExecCount := (monthsToEval * 3) //1, 2, 3 of each month
	if execsCount != expectedExecCount {
		t.Errorf("Expected %d executions (%d done) in %d days (%d ticks).", expectedExecCount, execsCount, daysMaxToEval, ticksCount)
		return
	}
}
