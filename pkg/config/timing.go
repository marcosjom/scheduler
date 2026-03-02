package config

import (
	"errors"
	"strconv"
	"strings"
)

// years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
type Age string

// month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
type Range string

type Timing struct {

	// Defines the range on which the action should triggered.
	Range struct {
		// Examples are:
		// For yearly tasks: min: "december 31", max "december 31"
		// For monthly tasks: min: "1", max "1"
		// For weekly tasks: min: "monday 00h 00m", max "monday 00h 59m"
		// For daily tasks: min: "00h 00m", max "00h 59m"
		Min Range //month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
		Max Range //month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
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

// Lexical validation.
func (t Age) HasError() error {
	partsCount := 0
	parts := strings.Split(string(t), " ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		//eval years+"y" | months+"M" | weeks+"w" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		if len(p) == 1 {
			return errors.New("Part requires value and suffix: '" + p + "'.")
		}
		sfx := p[len(p)-1]
		switch sfx {
		case 'y', 'M', 'w', 'd', 'h', 'm', 's':
			num, err := strconv.Atoi(p[:len(p)-1])
			if err != nil {
				return errors.New("Part requires a numeric value: '" + p + "'.")
			} else if num <= 0 {
				return errors.New("Part requires a positive non-zero value: '" + p + "'.")
			}
			partsCount++
		default:
			return errors.New("Invalid suffix '" + string(sfx) + "'.")
		}
	}
	//
	if partsCount <= 0 {
		return errors.New("Empty age's value.")
	}
	//
	return nil
}

// Lexical validation.
func (t Range) HasError() error {
	partsCount := 0
	parts := strings.Split(string(t), " ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		//eval month-name
		if strings.EqualFold(p, "January") || strings.EqualFold(p, "Jan") || strings.EqualFold(p, "Jan.") || strings.EqualFold(p, "JA") ||
			strings.EqualFold(p, "February") || strings.EqualFold(p, "Feb") || strings.EqualFold(p, "Feb.") || strings.EqualFold(p, "FE") ||
			strings.EqualFold(p, "March") || strings.EqualFold(p, "Mar") || strings.EqualFold(p, "Mar.") || strings.EqualFold(p, "MR") ||
			strings.EqualFold(p, "April") || strings.EqualFold(p, "Apr") || strings.EqualFold(p, "Apr.") || strings.EqualFold(p, "AP") ||
			strings.EqualFold(p, "May") || strings.EqualFold(p, "MY") ||
			strings.EqualFold(p, "June") || strings.EqualFold(p, "Jun") || strings.EqualFold(p, "Jun.") || strings.EqualFold(p, "JN") ||
			strings.EqualFold(p, "July") || strings.EqualFold(p, "Jul") || strings.EqualFold(p, "Jul.") || strings.EqualFold(p, "JL") ||
			strings.EqualFold(p, "August") || strings.EqualFold(p, "Aug") || strings.EqualFold(p, "Aug.") || strings.EqualFold(p, "AU") ||
			strings.EqualFold(p, "September") || strings.EqualFold(p, "Sep") || strings.EqualFold(p, "Sep.") || strings.EqualFold(p, "Sept") || strings.EqualFold(p, "Sept.") || strings.EqualFold(p, "SE") ||
			strings.EqualFold(p, "October") || strings.EqualFold(p, "Oct") || strings.EqualFold(p, "Oct.") || strings.EqualFold(p, "OC") ||
			strings.EqualFold(p, "November") || strings.EqualFold(p, "Nov") || strings.EqualFold(p, "Nov.") || strings.EqualFold(p, "NV") ||
			strings.EqualFold(p, "December") || strings.EqualFold(p, "Dec") || strings.EqualFold(p, "Dec.") || strings.EqualFold(p, "DE") {
			//is a month-name
			partsCount++
			continue
		}
		//eval weekday-name
		if strings.EqualFold(p, "Monday") || strings.EqualFold(p, "Mon") || strings.EqualFold(p, "Mon.") || strings.EqualFold(p, "Mo") || strings.EqualFold(p, "Mo.") ||
			strings.EqualFold(p, "Tuesday") || strings.EqualFold(p, "Tue") || strings.EqualFold(p, "Tue.") || strings.EqualFold(p, "Tu") || strings.EqualFold(p, "Tu.") ||
			strings.EqualFold(p, "Wednesday") || strings.EqualFold(p, "Wed") || strings.EqualFold(p, "Wed.") || strings.EqualFold(p, "We") || strings.EqualFold(p, "We.") ||
			strings.EqualFold(p, "Thursday") || strings.EqualFold(p, "Thu") || strings.EqualFold(p, "Thu.") || strings.EqualFold(p, "Th") || strings.EqualFold(p, "Th.") ||
			strings.EqualFold(p, "Friday") || strings.EqualFold(p, "Fri") || strings.EqualFold(p, "Fri.") || strings.EqualFold(p, "Fr") || strings.EqualFold(p, "Fr.") ||
			strings.EqualFold(p, "Saturday") || strings.EqualFold(p, "Sat") || strings.EqualFold(p, "Sat.") || strings.EqualFold(p, "Sa") || strings.EqualFold(p, "Sa.") ||
			strings.EqualFold(p, "Sunday") || strings.EqualFold(p, "Sun") || strings.EqualFold(p, "Sun.") || strings.EqualFold(p, "Su") || strings.EqualFold(p, "Su.") {
			//is a weekday-name
			partsCount++
			continue
		}
		//eval numeric-day
		num, err := strconv.Atoi(p)
		if err == nil {
			if num < 1 || num > 31 {
				return errors.New("Day of month should be between 1 and 31: '" + p + "'.")
			}
			partsCount++
			continue
		}
		//eval hour+"h", minute+"m", second+"s"
		if len(p) == 1 {
			return errors.New("Part requires value and suffix: '" + p + "'.")
		}
		sfx := p[len(p)-1]
		switch sfx {
		case 'h', 'm', 's':
			num, err := strconv.Atoi(p[:len(p)-1])
			if err != nil {
				return errors.New("Part requires a numeric value: '" + p + "'.")
			} else if num < 0 {
				return errors.New("Part requires a positive value: '" + p + "'.")
			}
			partsCount++
		default:
			return errors.New("Invalid suffix '" + string(sfx) + "'.")
		}
	}
	//
	if partsCount <= 0 {
		return errors.New("Empty range's value.")
	}
	//
	return nil
}
