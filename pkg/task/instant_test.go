package task

import (
	"testing"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/common"
	"github.com/marcosjom/sys-backups-automation/pkg/config"
)

// Testing Instant.String()
func Test_Instant_String(t *testing.T) {
	v0 := Instant{Month: common.January,
		Weekday:  common.Wednesday,
		Monthday: common.Monthday(5),
		Hour:     common.Hour(13),
		Minute:   common.Minute(33),
		Second:   common.Second(59),
	}
	v0Str := config.Instant(v0.String())
	//
	v1 := Instant{}
	if err := v1.SetFromConfig(v0Str); err != nil {
		t.Errorf("Instant -> config.Instant('%s') -> Instant returned error: %s.", string(v0Str), err.Error())
		return
	} else if v0 != v1 {
		t.Errorf("Instant -> config.Instant('%s') -> Instant('%s'); produced a different result.", string(v0Str), v1.String())
		return
	}
}

func Test_Instant_Yearly(t *testing.T) {
	instant := Instant{Month: common.July}
	//
	year := 2026
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	curTime := start
	ticksToEval := 12
	ticksShouldBeTrue := 6
	ticksWereTrue := 0
	for iTick := 0; iTick < ticksToEval; iTick++ {
		//eval
		//isTrue := "no"
		relLoc := instant.GetLocationFrom(curTime)
		if relLoc == InstantBefore || relLoc == InstantEqual {
			//isTrue = "yes"
			ticksWereTrue++
		}
		//log.Printf("curTime: %s = %s\n", curTime.String(), isTrue)
		//move to next month
		latDayInMonth := common.LastDayInMonth(year, int(curTime.Month()))
		curTime = curTime.Add(time.Hour * time.Duration(24*latDayInMonth))
	}
	if ticksShouldBeTrue != ticksWereTrue {
		t.Errorf("Expected %d ticks to be true, found %d (of %d).", ticksShouldBeTrue, ticksWereTrue, ticksToEval)
		return
	}
}

func Test_Instant_Monthly(t *testing.T) {
	instant := Instant{Monthday: common.Monthday(8)} //1-7
	year := 2026
	month := time.Month(1)
	start := time.Date(year, month, 1, 13, 31, 59, 0, time.UTC)
	curTime := start
	monthsToEval := 6
	monthsEvaluated := 0
	curMonth := month
	ticksShouldBeTrue := (monthsToEval * 7) //1..7
	ticksWereTrue := 0
	for monthsEvaluated < monthsToEval {
		//eval
		//isTrue := "no"
		relLoc := instant.GetLocationFrom(curTime)
		if relLoc == InstantAfter {
			//isTrue = "yes"
			ticksWereTrue++
		}
		//log.Printf("curTime: %s = %s\n", curTime.String(), isTrue)
		//move to next day
		curTime = curTime.Add(time.Hour * time.Duration(24))
		if curMonth != curTime.Month() {
			curMonth = curTime.Month()
			monthsEvaluated++
		}
	}
	if ticksShouldBeTrue != ticksWereTrue {
		t.Errorf("Expected %d ticks to be true, found %d (of %d months evaluated).", ticksShouldBeTrue, ticksWereTrue, monthsEvaluated)
		return
	}
}

func Test_Instant_Weekly(t *testing.T) {
	instant := Instant{Weekday: common.Thursday}
	//
	year := 2026
	start := time.Date(year, 1, 1, 13, 31, 59, 0, time.UTC)
	curTime := start
	weeksToEval := 6
	ticksToEval := (weeksToEval * 7)
	ticksShouldBeTrue := (weeksToEval * 3) //Thursday, Friday, Saturday
	ticksWereTrue := 0
	for iTick := 0; iTick < ticksToEval; iTick++ {
		//eval
		relLoc := instant.GetLocationFrom(curTime)
		if relLoc == InstantBefore || relLoc == InstantEqual {
			ticksWereTrue++
		}
		//move to next day
		curTime = curTime.Add(time.Hour * time.Duration(24))
	}
	if ticksShouldBeTrue != ticksWereTrue {
		t.Errorf("Expected %d ticks to be true, found %d (of %d).", ticksShouldBeTrue, ticksWereTrue, ticksToEval)
		return
	}
}

func Test_Instant_Daily(t *testing.T) {
	instant := Instant{Hour: common.Hour(17)}
	//
	start := time.Date(2026, 1, 1, 13, 31, 59, 0, time.UTC)
	curTime := start
	daysToEval := 31
	ticksToEval := (daysToEval * 24)
	ticksShouldBeTrue := (daysToEval * (24 - 16))
	ticksWereTrue := 0
	for iTick := 0; iTick < ticksToEval; iTick++ {
		//eval
		relLoc := instant.GetLocationFrom(curTime)
		if relLoc == InstantBefore || relLoc == InstantEqual {
			ticksWereTrue++
		}
		//move to next day
		curTime = curTime.Add(time.Hour * time.Duration(1))
	}
	if ticksShouldBeTrue != ticksWereTrue {
		t.Errorf("Expected %d ticks to be true, found %d (of %d).", ticksShouldBeTrue, ticksWereTrue, ticksToEval)
		return
	}
}
