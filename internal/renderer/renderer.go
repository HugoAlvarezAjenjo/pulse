package renderer

import (
	"fmt"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// Renderer defines the interface for outputting check results.
type Renderer interface {
	// Header prints the initial header before checks run.
	Header()
	// Success renders a passing check result.
	Success(r result.Result)
	// Failure renders a failing check result.
	Failure(r result.Result)
	// Error renders an internal error result.
	Error(r result.Result)
	// Summary renders the final summary after all checks.
	Summary(results []result.Result, duration time.Duration)
}

// Render outputs all results using the given renderer.
func Render(rnd Renderer, results []result.Result, duration time.Duration) {
	rnd.Header()

	for _, r := range results {
		switch r.Status {
		case result.Success:
			rnd.Success(r)
		case result.Failure:
			rnd.Failure(r)
		case result.Error:
			rnd.Error(r)
		}
	}

	rnd.Summary(results, duration)
}

// formatDuration formats a duration in a human-friendly way.
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
