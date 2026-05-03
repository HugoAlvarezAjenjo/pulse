package checks

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// EnvCheck validates that an environment variable exists and optionally
// matches an expected pattern.
type EnvCheck struct {
	Name     string
	Variable string
	Expected string
	Fix      *result.Fix
}

// Run checks if the environment variable is set.
func (c *EnvCheck) Run(_ context.Context) result.Result {
	start := time.Now()

	value, exists := os.LookupEnv(c.Variable)
	duration := time.Since(start)

	if !exists {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("environment variable %s is not set", c.Variable),
			Hint:     fmt.Sprintf("set %s in your environment or .env file", c.Variable),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	if value == "" {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("environment variable %s is empty", c.Variable),
			Hint:     fmt.Sprintf("set a value for %s", c.Variable),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	// Check expected pattern if provided
	if c.Expected != "" && !matchPattern(value, c.Expected) {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("%s does not match expected pattern %q", c.Variable, c.Expected),
			Hint:     fmt.Sprintf("expected %s to match %q", c.Variable, c.Expected),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  fmt.Sprintf("%s is set", c.Variable),
		Duration: duration,
		Fix:      c.Fix,
	}
}

// matchPattern performs simple wildcard matching.
// Supports * as a wildcard for any characters.
func matchPattern(value, pattern string) bool {
	// Simple prefix match: "sk-*"
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(value, prefix)
	}

	// Simple suffix match: "*-production"
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(value, suffix)
	}

	// Exact match
	return value == pattern
}
