package task

type Result int

const (
	// Task execution completed succesfuly.
	Success Result = iota
	// Task completion failed ot execute, but should retry later.
	ErrorRecoverable
	// Task completion failed ot execute, and should not retry until configuration changes.
	ErrorUnrecoverable
)

func (r *Result) AsText() string {
	switch *r {
	case Success:
		return "Success"
	case ErrorRecoverable:
		return "ErrorRecoverable"
	case ErrorUnrecoverable:
		return "ErrorUnrecoverable"
	}
	return "Unknown"
}
