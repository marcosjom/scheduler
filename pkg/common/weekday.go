package common

import (
	"errors"
	"strings"
	"time"
)

type Weekday uint8

const (
	Weekday_invalid = iota // 0
	Sunday
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Weekday) String() string {
	switch d {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	}
	return ""
}

func (d Weekday) Index() uint8 {
	switch d {
	case Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday:
		return uint8(d) - 1
	}
	return Weekday_invalid
}

func WeekdayFromStr(p string) (Weekday, error) {
	if strings.EqualFold(p, "Monday") || strings.EqualFold(p, "Mon") || strings.EqualFold(p, "Mon.") || strings.EqualFold(p, "Mo") || strings.EqualFold(p, "Mo.") {
		return Monday, nil
	}
	if strings.EqualFold(p, "Tuesday") || strings.EqualFold(p, "Tue") || strings.EqualFold(p, "Tue.") || strings.EqualFold(p, "Tu") || strings.EqualFold(p, "Tu.") {
		return Tuesday, nil
	}
	if strings.EqualFold(p, "Wednesday") || strings.EqualFold(p, "Wed") || strings.EqualFold(p, "Wed.") || strings.EqualFold(p, "We") || strings.EqualFold(p, "We.") {
		return Wednesday, nil
	}
	if strings.EqualFold(p, "Thursday") || strings.EqualFold(p, "Thu") || strings.EqualFold(p, "Thu.") || strings.EqualFold(p, "Th") || strings.EqualFold(p, "Th.") {
		return Thursday, nil
	}
	if strings.EqualFold(p, "Friday") || strings.EqualFold(p, "Fri") || strings.EqualFold(p, "Fri.") || strings.EqualFold(p, "Fr") || strings.EqualFold(p, "Fr.") {
		return Friday, nil
	}
	if strings.EqualFold(p, "Saturday") || strings.EqualFold(p, "Sat") || strings.EqualFold(p, "Sat.") || strings.EqualFold(p, "Sa") || strings.EqualFold(p, "Sa.") {
		return Saturday, nil
	}
	if strings.EqualFold(p, "Sunday") || strings.EqualFold(p, "Sun") || strings.EqualFold(p, "Sun.") || strings.EqualFold(p, "Su") || strings.EqualFold(p, "Su.") {
		return Sunday, nil
	}
	return Month_invalid, errors.New("Not a wekkday-string.")
}

func WeekdayFromTime(wd time.Time) Weekday {
	switch wd.Weekday() {
	case time.Sunday:
		return Sunday
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	default:
		return Weekday_invalid
	}
}
