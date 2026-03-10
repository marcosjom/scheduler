package task

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/common"
	"github.com/marcosjom/sys-backups-automation/pkg/config"
)

type AgePart int

const (
	Seconds AgePart = iota
	Minutes
	Hours
	Days
	Months
	Years
)

type Age struct {
	Seconds int
	Minutes int
	Hours   int
	Days    int
	Months  int
	Years   int
}

func (a *Age) String() string {
	r := ""
	if a.Years != 0 {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Years) + "y"
	}
	if a.Months != 0 {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Months) + "M"
	}
	if a.Days != 0 {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Days) + "d"
	}
	if a.Hours != 0 {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Hours) + "h"
	}
	if a.Minutes != 0 {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Minutes) + "m"
	}
	if a.Seconds != 0 || r == "" {
		if r != "" {
			r += " "
		}
		r += strconv.Itoa(a.Seconds) + "s"
	}
	return r
}

func (a *Age) AsSecs() int64 {
	return int64(a.Seconds) +
		int64(a.Minutes*60) +
		int64(a.Hours*60*60) +
		int64(a.Days*60*60*24) +
		int64(a.Months*60*60*24*31) +
		int64(a.Years*60*60*24*31*12)
}

func (a *Age) SetFromConfig(age config.Age) error {
	tmp := Age{}
	partsCount := 0
	parts := strings.Split(string(age), " ")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		//eval years+"y" | months+"M" | days+"d" | hours+"h" | minutes+"m" | seconds+"s"
		if len(p) == 1 {
			return errors.New("Part requires value and suffix: '" + p + "'.")
		}
		sfx := p[len(p)-1]
		switch sfx {
		case 'y', 'M', 'd', 'h', 'm', 's':
			num, err := strconv.Atoi(p[:len(p)-1])
			if err != nil {
				return errors.New("Part requires a numeric value: '" + p + "'.")
			} else if num <= 0 {
				return errors.New("Part requires a positive non-zero value: '" + p + "'.")
			}
			switch sfx {
			case 'y':
				if tmp.Years != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Years = num
			case 'M':
				if tmp.Months != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Months = num
			case 'd':
				if tmp.Days != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Days = num
			case 'h':
				if tmp.Hours != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Hours = num
			case 'm':
				if tmp.Minutes != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Minutes = num
			case 's':
				if tmp.Seconds != 0 {
					return errors.New("Part '" + string(sfx) + "' appears more than once.")
				}
				tmp.Seconds = num
			default:
				return errors.New("Iternal error with '" + string(sfx) + "'.")
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
	// Set
	*a = tmp
	return nil
}

func (a *Age) SetFromTo(startP time.Time, endP time.Time) {
	tmp := Age{}
	start := startP
	end := endP
	// Force Start <= End
	if end.Sub(start) < 0 {
		end = startP
		start = endP
	}
	// Parts (months and days are indices too)
	startYear, startMonth, startDay, startHour, startMin, startSec := start.Year(), int(start.Month())-1, int(start.Day())-1, start.Hour(), start.Minute(), start.Second()
	endYear, endMonth, endDay, endHour, endMin, endSec := end.Year(), int(end.Month())-1, int(end.Day())-1, end.Hour(), end.Minute(), end.Second()
	//years/months
	{
		// Calculate possible-completed months
		sMonths := (startYear * 12) + startMonth
		eMonths := (endYear * 12) + endMonth
		monthsCount := (eMonths - sMonths)
		// Remove incomplete month
		if endDay < startDay ||
			(endDay == startDay && endHour < startHour) ||
			(endDay == startDay && endHour == startHour && endMin < startMin) ||
			(endDay == startDay && endHour == startHour && endMin == startMin && endSec < startSec) {
			monthsCount -= 1
			lastDayInMonth := common.LastDayInMonth(endYear, startMonth+1) //start-month at end-year
			endDay = lastDayInMonth + endDay
		}
		// Set years and months
		if monthsCount > 0 {
			tmp.Years = monthsCount / 12
			tmp.Months = monthsCount % 12
		}
	}
	//days/hours
	{
		// Force sDay <= eDay2
		eDay2 := endDay
		if eDay2 < startDay {
			lastDayInMonth := common.LastDayInMonth(endYear, startMonth+1) //start-month at end-year
			eDay2 = lastDayInMonth + eDay2
		}
		// Calculate possible-completed hours
		sHours := (startDay * 24) + startHour
		eHours := (eDay2 * 24) + endHour
		hoursCount := (eHours - sHours)
		// Remove incomplete hour
		if endMin < startMin ||
			(endMin == startMin && endSec < startSec) {
			hoursCount -= 1
			endMin += 60
		}
		// Set days and hours
		if hoursCount > 0 {
			tmp.Days = hoursCount / 24
			tmp.Hours = hoursCount % 24
		}
	}
	//minutes/seconds
	{
		// Force sMinute <= eMinute2
		eMin2 := endMin
		if eMin2 < startMin {
			eMin2 = 60 + eMin2
		}
		// Calculate seconds
		sSecs := (startMin * 60) + startSec
		eSecs := (eMin2 * 60) + endSec
		secsCount := (eSecs - sSecs)
		// Set mins and secs
		tmp.Minutes = (secsCount / 60)
		tmp.Seconds = (secsCount % 60)
	}
	// Set
	*a = tmp
}
