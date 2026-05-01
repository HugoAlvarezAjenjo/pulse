package checks

import (
	"context"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// Check defines the interface for all environment checks.
type Check interface {
	// Run executes the check and returns a result.
	Run(ctx context.Context) result.Result
}
