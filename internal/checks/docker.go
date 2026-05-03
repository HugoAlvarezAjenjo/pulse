package checks

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// DockerCheck validates that a Docker container is running.
type DockerCheck struct {
	Name      string
	Container string
	Fix       *result.Fix
}

// Run checks if the specified Docker container is running.
func (c *DockerCheck) Run(ctx context.Context) result.Result {
	start := time.Now()

	// First check if docker CLI is available
	if _, err := dockerLookPath(); err != nil {
		return result.Result{
			Name:     c.Name,
			Status:   result.Error,
			Message:  "docker CLI not found",
			Hint:     "install Docker: https://docs.docker.com/get-docker/",
			Duration: time.Since(start),
			Fix:      c.Fix,
		}
	}

	// Check container state using docker inspect
	cmd := exec.CommandContext(ctx, "docker", "inspect",
		"--format", "{{.State.Running}}", c.Container)

	output, err := cmd.Output()
	duration := time.Since(start)

	if err != nil {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("container %q not found or not accessible", c.Container),
			Hint:     fmt.Sprintf("start the container: docker start %s", c.Container),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	running := strings.TrimSpace(string(output))
	if running != "true" {
		return result.Result{
			Name:     c.Name,
			Status:   result.Failure,
			Message:  fmt.Sprintf("container %q exists but is not running", c.Container),
			Hint:     fmt.Sprintf("start the container: docker start %s", c.Container),
			Duration: duration,
			Fix:      c.Fix,
		}
	}

	return result.Result{
		Name:     c.Name,
		Status:   result.Success,
		Message:  fmt.Sprintf("container %q is running", c.Container),
		Duration: duration,
		Fix:      c.Fix,
	}
}

// dockerLookPath finds the docker executable in PATH.
func dockerLookPath() (string, error) {
	name := "docker"
	if runtime.GOOS == "windows" {
		name = "docker.exe"
	}
	return exec.LookPath(name)
}
