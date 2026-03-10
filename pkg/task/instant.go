package task

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/common"
	"github.com/marcosjom/sys-backups-automation/pkg/config"
)

type InstantPart int

const (
	//1-12
	Month InstantPart = iota
	//1-7
	Weekday
	//1-31
	Monthday
	//1-25
	Hour
	//1-61
	Minute
	//1-61
	Second
)

type InstantLocation int

const (
	InstantUnknown InstantLocation = iota
	InstantBefore
	InstantEqual
	InstantAfter
)

type Instant struct {
	//1-12
	Month common.Month
	//1-7
	Weekday common.Weekday
	//1-31
	Monthday common.Monthday
	//1-25
	Hour common.Hour
	//1-61
	Minute common.Minute
	//1-61
	Second common.Second
}

func (i *Instant) String() string {
	r := ""
	if i.Month > common.Month_invalid {
		if r != "" {
			r += " "
		}
		r += i.Month.String()
	}
	if i.Weekday > common.Weekday_invalid {
		if r != "" {
			r += " "
		}
		r += i.Weekday.String()
	}
	if i.Monthday > common.Monthday_invalid {
		if r != "" {
			r += " "
		}
		r += i.Monthday.String()
	}
	if i.Hour > common.Hour_invalid {
		if r != "" {
			r += " "
		}
		r += i.Hour.String() + "h"
	}
	if i.Minute > common.Minute_invalid {
		if r != "" {
			r += " "
		}
		r += i.Minute.String() + "m"
	}
	if i.Second > common.Second_invalid {
		if r != "" {
			r += " "
		}
		r += i.Second.String() + "s"
	}
	return r
}

func (i *Instant) SetFromConfig(instant config.Instant) error {
	tmp := Instant{}
	partsCount := 0
	parts := strings.Split(string(instant), " ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		//eval month-name
		if m, err := common.MonthFromStr(p); err == nil {
			//is a month-name
			if tmp.Month != common.Month_invalid {
				return errors.New("Instant has more than one month component: '" + string(instant) + "'.")
			}
			tmp.Month = m
			partsCount++
			continue
		}
		//eval weekday-name
		if wd, err := common.WeekdayFromStr(p); err == nil {
			//is a weekday-name
			if tmp.Weekday != common.Weekday_invalid {
				return errors.New("Instant has more than one weekday component: '" + string(instant) + "'.")
			}
			tmp.Weekday = wd
			partsCount++
			continue
		}
		//eval numeric-day
		if num, err := strconv.Atoi(p); err == nil {
			if num < 1 || num > 31 {
				return errors.New("Day of month should be between 1 and 31: '" + p + "'.")
			}
			if tmp.Monthday != common.Month_invalid {
				return errors.New("Instant has more than one monthday component: '" + string(instant) + "'.")
			}
			tmp.Monthday = common.Monthday(num)
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
			// Apply number
			switch sfx {
			case 'h':
				if tmp.Hour != common.Hour_invalid {
					return errors.New("Instant has more than one hour component: '" + string(instant) + "'.")
				}
				tmp.Hour = common.Hour(num + 1)
			case 'm':
				if tmp.Minute != common.Minute_invalid {
					return errors.New("Instant has more than one minute component: '" + string(instant) + "'.")
				}
				tmp.Minute = common.Minute(num + 1)
			case 's':
				if tmp.Second != common.Second_invalid {
					return errors.New("Instant has more than one second component: '" + string(instant) + "'.")
				}
				tmp.Second = common.Second(num + 1)
			default:
				return errors.New("Internal error, unexpected suffix('" + string(sfx) + "') after validation.")
			}
			partsCount++
		default:
			return errors.New("Invalid suffix '" + string(sfx) + "'.")
		}
	}
	//
	if partsCount <= 0 {
		return errors.New("Empty instant's value.")
	}
	// Set
	*i = tmp
	return nil
}

func (i *Instant) GetLocationFrom(time time.Time) InstantLocation {
	isActivated := false
	// Compare Month
	if i.Month > common.Month_invalid {
		v0 := i.Month
		v1 := common.Month(time.Month())
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same month, keep evaluating
		isActivated = true
	}
	// Compare Monthday
	if isActivated || i.Monthday > common.Monthday_invalid {
		isActivated = true
		v1 := common.Monthday(time.Day())
		v0 := v1
		if i.Monthday > common.Monthday_invalid {
			v0 = i.Monthday
		}
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same monthday, keep evaluating
		isActivated = true
	}
	// Compare Weekday
	if isActivated || i.Weekday > common.Weekday_invalid {
		isActivated = true
		v1 := common.WeekdayFromTime(time)
		v0 := v1
		if i.Weekday > common.Weekday_invalid {
			v0 = i.Weekday
		}
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same weekday, keep evaluating
		isActivated = true
	}
	// Compare Hour
	if isActivated || i.Hour > common.Hour_invalid {
		isActivated = true
		v1 := common.Hour(time.Hour() + 1)
		v0 := v1
		if i.Hour > common.Hour_invalid {
			v0 = i.Hour
		}
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same hour, keep evaluating
		isActivated = true
	}
	// Compare Minute
	if isActivated || i.Minute > common.Minute_invalid {
		isActivated = true
		v1 := common.Minute(time.Minute() + 1)
		v0 := v1
		if i.Minute > common.Minute_invalid {
			v0 = i.Minute
		}
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same minute, keep evaluating
		isActivated = true
	}
	// Compare Second
	if isActivated || i.Second > common.Second_invalid {
		isActivated = true
		v1 := common.Second(time.Second() + 1)
		v0 := v1
		if i.Second > common.Second_invalid {
			v0 = i.Second
		}
		if v0 < v1 {
			return InstantBefore
		} else if v0 > v1 {
			return InstantAfter
		}
		// Same second, is equal
		return InstantEqual
	}
	//nothing to compare
	return InstantUnknown
}
