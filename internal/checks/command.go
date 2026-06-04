package checks

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// CommandCheck validates that a command exists and optionally
// matches a version constraint.
type CommandCheck struct {
	Name     string
	Command  string
	Expected string
	Fix      *result.Fix
}

// Run executes the command and validates the output.
func (c *CommandCheck) Run(ctx context.Context) result.Result {
	start := time.Now()

	command := strings.TrimSpace(c.Command)
	if command == "" {
		return result.Result{
			Name:     c.Name,
			Status:   result.Error,
			Message:  "empty command",
			Duration: time.Since(start),
			Fix:      c.Fix,
		}
	}

	cmd := shellCommandContext(ctx, command)
	output, err := cmd.Output()
	duration := time.Since(start)

	if err != nil {
		parts := strings.Fields(command)
		hint := "ensure the command is valid and available in PATH"
		if len(parts) > 0 {
			hint = fmt.Sprintf("ensure '%s' is installed and in PATH", parts[0])
		}
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("command failed: %s", err),
			Hint:     hint,
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	version := strings.TrimSpace(string(output))

	// If no expected constraint, just check command exists
	if c.Expected == "" {
		return result.Result{
			Name:     c.Name,
			Status:   result.Success,
			Message:  version,
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	// Parse and validate version constraint
	if !matchesConstraint(version, c.Expected) {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("version %s does not satisfy %s", version, c.Expected),
			Hint:     fmt.Sprintf("expected %s", c.Expected),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  version,
		Duration: duration,
		Fix:      c.Fix,
	}
}

// matchesConstraint checks if the version string satisfies the constraint.
// Supports semver constraints via hashicorp/go-version (>=, <=, >, <, ~>, !=, and ranges).
func matchesConstraint(version, constraint string) bool {
	raw := extractVersion(version)
	v, err := goversion.NewVersion(raw)
	if err != nil {
		return strings.Contains(version, constraint)
	}

	c, err := goversion.NewConstraint(constraint)
	if err != nil {
		return strings.Contains(version, constraint)
	}

	return c.Check(v)
}

var versionRegex = regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)

func extractVersion(output string) string {
	match := versionRegex.FindString(output)
	if match != "" {
		return match
	}
	return strings.TrimSpace(output)
}

func shellCommandContext(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}
