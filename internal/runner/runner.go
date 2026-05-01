package runner

import (
	"context"
	"sync"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/checks"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

const defaultTimeout = 30 * time.Second

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
func (r *Runner) Run(ctx context.Context, checkList []checks.Check) []result.Result {
	results := make([]result.Result, len(checkList))
	var wg sync.WaitGroup

	timeout := r.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	for i, check := range checkList {
		wg.Add(1)
		go func(idx int, c checks.Check) {
			defer wg.Done()

			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			results[idx] = c.Run(checkCtx)
		}(i, check)
	}

	wg.Wait()
	return results
}
