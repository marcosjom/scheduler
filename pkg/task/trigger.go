package task

import (
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/config"
)

// Determines if a task should be executed.
type Trigger struct {
	History *History
	// Cache
	cache struct {
		config struct {
			crc32   uint32
			isValid bool
			// The timing for triggering this task.
			timing struct {
				rangee struct {
					min Instant
					max Instant
				}
				age struct {
					tick     Age
					min      Age
					maxIsSet bool
					max      Age
				}
			}
		}
	}
}

// Config

func (t *Trigger) SetConfig(config *config.Task) error {
	//skip evaluating same config multiple times
	crc32 := config.Crc32()
	if t.cache.config.crc32 == crc32 {
		return nil
	}
	//header
	t.cache.config.crc32 = crc32
	t.cache.config.isValid = true
	//timing
	if err := t.cache.config.timing.rangee.min.SetFromConfig(config.Timing.Range.Min); err != nil {
		// Config error (dont try again until fixed)
		t.cache.config.isValid = false
		return err
	}
	if err := t.cache.config.timing.rangee.max.SetFromConfig(config.Timing.Range.Max); err != nil {
		// Config error (dont try again until fixed)
		t.cache.config.isValid = false
		return err
	}
	//age
	if err := t.cache.config.timing.age.tick.SetFromConfig(config.Timing.Age.Tick); err != nil {
		// Config error (dont try again until fixed)
		t.cache.config.isValid = false
		return err
	}
	if err := t.cache.config.timing.age.min.SetFromConfig(config.Timing.Age.Min); err != nil {
		// Config error (dont try again until fixed)
		t.cache.config.isValid = false
		return err
	}
	t.cache.config.timing.age.maxIsSet = true
	if err := t.cache.config.timing.age.max.SetFromConfig(config.Timing.Age.Max); err != nil {
		t.cache.config.timing.age.maxIsSet = false
	}
	//
	return nil
}

// Tick's age

func (t *Trigger) IsTickOldEnoughNow() bool {
	return t.IsTickOldEnough(time.Now())
}

func (t *Trigger) IsTickOldEnough(time time.Time) bool {

	// Validate Config
	if !t.cache.config.isValid {
		return false
	}

	// Never attempted?
	if t.History.LastTick.Time.IsZero() {
		return true
	}

	// Determine current age
	age := Age{}
	age.SetFromTo(t.History.LastTick.Time, time)
	secsMin := t.cache.config.timing.age.tick.AsSecs()
	secsAge := age.AsSecs()

	// Compare ages
	if secsAge < secsMin {
		return false
	}

	// Is Old enough
	return true
}

// Success's age

func (t *Trigger) IsOldEnoughNow() bool {
	return t.IsOldEnough(time.Now())
}

func (t *Trigger) IsOldEnough(time time.Time) bool {

	// Validate Config
	if !t.cache.config.isValid {
		return false
	}

	// Never executed?
	if t.History.LastSuccess.Time.IsZero() {
		return true
	}

	// Determine current age
	age := Age{}
	age.SetFromTo(t.History.LastSuccess.Time, time)
	secsMin := t.cache.config.timing.age.min.AsSecs()
	secsAge := age.AsSecs()

	// Compare ages
	if secsAge < secsMin {
		return false
	}

	// Is Old enought
	return true
}

// Unskippable success's age

func (t *Trigger) IsUnacceptableOldNow() bool {
	return t.IsUnacceptableOld(time.Now())
}

func (t *Trigger) IsUnacceptableOld(time time.Time) bool {

	// Validate Config
	if !t.cache.config.isValid {
		return false
	}

	// Never executed?
	if t.History.LastSuccess.Time.IsZero() {
		return false
	}

	// Never unaceptable old (by config)
	if !t.cache.config.timing.age.maxIsSet {
		return false
	}

	// Determine current age
	age := Age{}
	age.SetFromTo(t.History.LastSuccess.Time, time)
	secsMax := t.cache.config.timing.age.max.AsSecs()
	secsAge := age.AsSecs()

	// Compare ages
	if secsAge < secsMax {
		return false
	}

	// Is unaceptable old
	return true
}

//

func (t *Trigger) ShouldRunTaskNow() bool {
	return t.ShouldRunTask(time.Now())
}

func (t *Trigger) ShouldRunTask(time time.Time) bool {

	// Validate Config
	if !t.cache.config.isValid {
		return false
	}

	// If last execution was non-recoverable-error;
	// do not try again unless configuration changed.
	hasCfgChanged := (t.cache.config.crc32 != t.History.LastRun.ConfigCrc32)
	if (t.History.LastRun.Result == ErrorUnrecoverable) && !hasCfgChanged {
		return false
	}

	// Determine last-success's age
	if !t.IsOldEnough(time) {
		return false
	}

	// Eval unacceptable old (overrides the range of execution)
	if t.IsUnacceptableOld(time) {
		return true
	}

	// Eval range of execution
	minLocation := t.cache.config.timing.rangee.min.GetLocationFrom(time)
	maxLocation := t.cache.config.timing.rangee.max.GetLocationFrom(time)
	if (minLocation == InstantBefore || minLocation == InstantEqual) &&
		(maxLocation == InstantAfter || maxLocation == InstantEqual) {
		return true
	}

	return false
}

func (t *Trigger) Execute(time time.Time) error {
	return nil
}
