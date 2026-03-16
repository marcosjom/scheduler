package task

import (
	"testing"
	"time"

	"github.com/marcosjom/scheduler/pkg/config"
)

// Testing Age.String()
func Test_Age_String(t *testing.T) {
	v0 := Age{Years: 2,
		Months:  3,
		Days:    4,
		Hours:   5,
		Minutes: 6,
		Seconds: 7,
	}
	v0Str := config.Age(v0.String())
	//
	v1 := Age{}
	if err := v1.SetFromConfig(v0Str); err != nil {
		t.Errorf("Age -> config.Age('%s') -> Age; produced error: %s", string(v0Str), err.Error())
		return
	} else if v0 != v1 {
		t.Errorf("Age -> config.Age('%s') -> Age('%s'); produced a different result.", string(v0Str), v1.String())
		return
	}
}

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
		newAgeSecs := newAge.AsSecs()
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
		newAgeSecs := newAge.AsSecs()
		if prevAgeSecs >= newAgeSecs {
			t.Errorf("Expected incremental age '%s' -> '%s'; ('%s' -> '%s')", start.String()[:19], end.String()[:19], prevAge.String(), newAge.String())
			return
		}
		prevAge = newAge
		prevAgeSecs = newAgeSecs
	}
}
