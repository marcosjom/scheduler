package config

import (
	"errors"
	"strconv"
	"strings"
)

// Instant represents a repeatable moment in time.
// month-name, weekday-name, numeric-day, hour+"h", minute+"m", second+"s"
type Instant string

// Lexical validation.
func (t Instant) HasError() error {
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
