package checks

import (
	"context"
	"testing"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestDockerCheck_NotFound(t *testing.T) {
	// This test assumes docker is installed but the container doesn't exist
	if _, err := dockerLookPath(); err != nil {
		t.Skip("docker not installed, skipping")
	}

	check := &DockerCheck{
		Name:      "nonexistent container",
		Container: "pulse_test_nonexistent_container_xyz",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := check.Run(ctx)
	if r.Status != result.Failure {
		t.Errorf("expected failure for nonexistent container, got %s: %s", r.Status, r.Message)
	}
}

func TestDockerCheck_NoDocker(t *testing.T) {
	// This test only makes sense if docker is NOT installed
	if _, err := dockerLookPath(); err == nil {
		t.Skip("docker is installed, skipping no-docker test")
	}

	check := &DockerCheck{
		Name:      "test container",
		Container: "redis",
	}

	r := check.Run(context.Background())
	if r.Status != result.Error {
		t.Errorf("expected error when docker not found, got %s", r.Status)
	}
}
