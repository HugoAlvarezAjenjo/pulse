package result

import "time"

// Status represents the outcome of a check.
type Status int

const (
	// Success indicates the check passed.
	Success Status = iota
	// Failure indicates the check did not pass.
	Failure
	// Error indicates an internal/runtime error during execution.
	Error
)

// String returns a human-readable representation of the status.
func (s Status) String() string {
	switch s {
	case Success:
		return "success"
	case Failure:
		return "failure"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Fix represents a user-defined remediation command.
type Fix struct {
	Run string
}

// Result holds the outcome of a single check execution.
type Result struct {
	Name     string
	Status   Status
	Message  string
	Hint     string
	Duration time.Duration
	Fix      *Fix
}

// Passed returns true if the check succeeded.
func (r Result) Passed() bool {
	return r.Status == Success
}
