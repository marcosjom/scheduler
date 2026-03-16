package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/marcosjom/scheduler/pkg/common"
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
		if _, err := common.MonthFromStr(p); err == nil {
			//is a month-name
			partsCount++
			continue
		}
		//eval weekday-name
		if _, err := common.WeekdayFromStr(p); err == nil {
			//is a weekday-name
			partsCount++
			continue
		}
		//eval numeric-day
		if num, err := strconv.Atoi(p); err == nil {
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
			} else if num > 59 {
				return errors.New("Part requires an under-60 value: '" + p + "'.")
			} else if sfx == 'h' && num > 23 {
				return errors.New("Part requires an under-24 value: '" + p + "'.")
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
