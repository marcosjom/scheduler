package task

import (
	"testing"
	"time"
)

// Testing calculating incremental age, one second increments during one year.
func Test_Age_FromTo_Calc_1Year_Incrementing_1Sec(t *testing.T) {
	start := time.Now()
	secsMaxToEval := 366 * 24 * 60 * 60
	prevAge := Age{}
	prevAgeSecs := int64(-1)
	for iSec := 0; iSec <= secsMaxToEval; iSec++ {
		end := start.Add(time.Second * time.Duration(iSec))
		newAge := Age{}
		newAge.SetFromTo(start, end)
		newAgeSecs := int64(newAge.Seconds) +
			int64(newAge.Minutes*60) +
			int64(newAge.Hours*60*60) +
			int64(newAge.Days*60*60*24) +
			int64(newAge.Months*60*60*24*31) +
			int64(newAge.Years*60*60*24*31*12)
		if prevAgeSecs >= newAgeSecs {
			t.Errorf("Expected incremental age '%s' -> '%s'; ('%s' -> '%s')", start.String()[:19], end.String()[:19], prevAge.String(), newAge.String())
			return
		}
		prevAge = newAge
		prevAgeSecs = newAgeSecs
	}
}

// Testing calculating incremental age, one hour increments during one year.
func Test_Age_FromTo_Calc_1Year_Incrementing_1Hour(t *testing.T) {
	start := time.Now()
	hoursMaxToEval := 366 * 24
	prevAge := Age{}
	prevAgeSecs := int64(-1)
	for iHour := 0; iHour <= hoursMaxToEval; iHour++ {
		end := start.Add(time.Hour * time.Duration(iHour))
		newAge := Age{}
		newAge.SetFromTo(start, end)
		newAgeSecs := int64(newAge.Seconds) +
			int64(newAge.Minutes*60) +
			int64(newAge.Hours*60*60) +
			int64(newAge.Days*60*60*24) +
			int64(newAge.Months*60*60*24*31) +
			int64(newAge.Years*60*60*24*31*12)
		if prevAgeSecs >= newAgeSecs {
			t.Errorf("Expected incremental age '%s' -> '%s'; ('%s' -> '%s')", start.String()[:19], end.String()[:19], prevAge.String(), newAge.String())
			return
		}
		prevAge = newAge
		prevAgeSecs = newAgeSecs
	}
}
