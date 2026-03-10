package common

import (
	"errors"
	"strings"
)

type Month uint8

const (
	Month_invalid = iota // 0
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (d Month) String() string {
	switch d {
	case January:
		return "January"
	case February:
		return "February"
	case March:
		return "March"
	case April:
		return "April"
	case May:
		return "May"
	case June:
		return "June"
	case July:
		return "July"
	case August:
		return "August"
	case September:
		return "September"
	case October:
		return "October"
	case November:
		return "November"
	case December:
		return "December"
	}
	return ""
}

func (d Month) Index() uint8 {
	switch d {
	case January, February, March, April, May, June, July, August, September, October, November, December:
		return uint8(d) - 1
	}
	return Month_invalid
}

func MonthFromStr(p string) (Month, error) {
	if strings.EqualFold(p, "January") || strings.EqualFold(p, "Jan") || strings.EqualFold(p, "Jan.") || strings.EqualFold(p, "JA") {
		return January, nil
	}
	if strings.EqualFold(p, "February") || strings.EqualFold(p, "Feb") || strings.EqualFold(p, "Feb.") || strings.EqualFold(p, "FE") {
		return February, nil
	}
	if strings.EqualFold(p, "March") || strings.EqualFold(p, "Mar") || strings.EqualFold(p, "Mar.") || strings.EqualFold(p, "MR") {
		return March, nil
	}
	if strings.EqualFold(p, "April") || strings.EqualFold(p, "Apr") || strings.EqualFold(p, "Apr.") || strings.EqualFold(p, "AP") {
		return April, nil
	}
	if strings.EqualFold(p, "May") || strings.EqualFold(p, "MY") {
		return May, nil
	}
	if strings.EqualFold(p, "June") || strings.EqualFold(p, "Jun") || strings.EqualFold(p, "Jun.") || strings.EqualFold(p, "JN") {
		return June, nil
	}
	if strings.EqualFold(p, "July") || strings.EqualFold(p, "Jul") || strings.EqualFold(p, "Jul.") || strings.EqualFold(p, "JL") {
		return July, nil
	}
	if strings.EqualFold(p, "August") || strings.EqualFold(p, "Aug") || strings.EqualFold(p, "Aug.") || strings.EqualFold(p, "AU") {
		return August, nil
	}
	if strings.EqualFold(p, "September") || strings.EqualFold(p, "Sep") || strings.EqualFold(p, "Sep.") || strings.EqualFold(p, "Sept") || strings.EqualFold(p, "Sept.") || strings.EqualFold(p, "SE") {
		return September, nil
	}
	if strings.EqualFold(p, "October") || strings.EqualFold(p, "Oct") || strings.EqualFold(p, "Oct.") || strings.EqualFold(p, "OC") {
		return October, nil
	}
	if strings.EqualFold(p, "November") || strings.EqualFold(p, "Nov") || strings.EqualFold(p, "Nov.") || strings.EqualFold(p, "NV") {
		return November, nil
	}
	if strings.EqualFold(p, "December") || strings.EqualFold(p, "Dec") || strings.EqualFold(p, "Dec.") || strings.EqualFold(p, "DE") {
		return December, nil
	}
	return Month_invalid, errors.New("Not a month-string.")
}

func LastDayInMonth(year int, month int) int {
	if month != 2 {
		return 31 - (month-1)%7%2
	}
	//feb
	if year&3 == 0 && (year%25 != 0 || year&15 == 0) {
		return 29
	}
	return 28
}
