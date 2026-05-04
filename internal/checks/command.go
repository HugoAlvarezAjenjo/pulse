package checks

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

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
// Supports >= prefix for minimum version comparison.
func matchesConstraint(version, constraint string) bool {
	if strings.HasPrefix(constraint, ">=") {
		minVersion := strings.TrimPrefix(constraint, ">=")
		return compareVersions(extractVersion(version), minVersion) >= 0
	}
	return strings.Contains(version, constraint)
}

// versionRegex matches the first semver-like pattern (e.g., "1.26.2", "22.3.0")
var versionRegex = regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)

// extractVersion finds the first version-like pattern in the output string.
// Handles any format: "v1.2.3", "go version go1.26.2 darwin/arm64", "node v22.3.0", etc.
func extractVersion(output string) string {
	match := versionRegex.FindString(output)
	if match != "" {
		return match
	}
	// Fallback: return trimmed input
	return strings.TrimSpace(output)
}

// compareVersions performs a simple semver comparison.
// Returns -1, 0, or 1.
func compareVersions(a, b string) int {
	aParts := parseVersion(a)
	bParts := parseVersion(b)

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			aNum = aParts[i]
		}
		if i < len(bParts) {
			bNum = bParts[i]
		}

		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}
	return 0
}

func shellCommandContext(ctx context.Context, command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.CommandContext(ctx, "cmd", "/c", command)
	}
	return exec.CommandContext(ctx, "sh", "-c", command)
}

// parseVersion splits a version string into numeric parts.
func parseVersion(v string) []int {
	parts := strings.Split(v, ".")
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		// Strip any non-numeric suffix (e.g., "22-beta")
		p = strings.Split(p, "-")[0]
		p = strings.Split(p, "+")[0]
		n, err := strconv.Atoi(p)
		if err != nil {
			nums = append(nums, 0)
		} else {
			nums = append(nums, n)
		}
	}
	return nums
}
