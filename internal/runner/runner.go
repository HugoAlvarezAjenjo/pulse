package runner

import (
	"context"
	"sync"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/checks"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

const defaultTimeout = 30 * time.Second

// CheckWithTimeout pairs a check with an optional per-check timeout.
type CheckWithTimeout struct {
	Check   checks.Check
	Timeout time.Duration // 0 means use global timeout
}

// Runner executes checks concurrently while preserving order.
type Runner struct {
	Timeout time.Duration
}

// New creates a new Runner with default settings.
func New() *Runner {
	return &Runner{
		Timeout: defaultTimeout,
	}
}

// Run executes all checks concurrently and returns results in the original order.
// Uses the global timeout for all checks.
func (r *Runner) Run(ctx context.Context, checkList []checks.Check) []result.Result {
	items := make([]CheckWithTimeout, len(checkList))
	for i, c := range checkList {
		items[i] = CheckWithTimeout{Check: c}
	}
	return r.RunWithTimeouts(ctx, items)
}

// RunWithTimeouts executes checks concurrently, respecting per-check timeouts.
// Priority: per-check timeout > global timeout > default (30s).
func (r *Runner) RunWithTimeouts(ctx context.Context, checkList []CheckWithTimeout) []result.Result {
	results := make([]result.Result, len(checkList))
	var wg sync.WaitGroup

	globalTimeout := r.Timeout
	if globalTimeout == 0 {
		globalTimeout = defaultTimeout
	}

	for i, item := range checkList {
		wg.Add(1)
		go func(idx int, cwt CheckWithTimeout) {
			defer wg.Done()

			// Per-check timeout takes priority over global
			timeout := globalTimeout
			if cwt.Timeout > 0 {
				timeout = cwt.Timeout
			}

			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			results[idx] = cwt.Check.Run(checkCtx)
		}(i, item)
	}

	wg.Wait()
	return results
}
