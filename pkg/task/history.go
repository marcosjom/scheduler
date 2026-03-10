package task

import (
	"time"
)

// Represents the execution history of a task.
type History struct {
	// Record's version (when a saving value changes, this value should be increased)
	VersionId uint32
	// Last Tick
	LastTick struct {
		Time time.Time
	}
	// Last Run
	LastRun struct {
		// Configuration that produced the last execution result.
		ConfigCrc32 uint32
		Time        time.Time
		Result      Result
	}
	// Last Success
	LastSuccess struct {
		Time time.Time
	}
	// Last Save (detects unsaved changes)
	LastSave struct {
		VersonId uint32
	}
}

func (h *History) SetLastRunResult(result Result) {
	h.LastRun.Time = time.Now()
	h.LastRun.Result = result
	//
	h.VersionId++
}

func (h *History) HasUnsavedChanges() bool {
	return h.LastSave.VersonId != h.VersionId
}

func (h *History) SyncAsSaved() {
	h.LastSave.VersonId = h.VersionId
}
