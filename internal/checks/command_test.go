package checks

import (
	"context"
	"testing"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestCommandCheck_Success(t *testing.T) {
	check := &CommandCheck{
		Name:    "echo test",
		Command: "echo hello",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
	if r.Message != "hello" {
		t.Errorf("expected message 'hello', got %q", r.Message)
	}
}

func TestCommandCheck_CommandNotFound(t *testing.T) {
	check := &CommandCheck{
		Name:    "nonexistent",
		Command: "nonexistent_command_xyz_123",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestCommandCheck_EmptyCommand(t *testing.T) {
	check := &CommandCheck{
		Name:    "empty",
		Command: "",
	}

	r := check.Run(context.Background())
	if r.Status != result.Error {
		t.Errorf("expected error, got %s", r.Status)
	}
}

func TestCommandCheck_VersionConstraint(t *testing.T) {
	check := &CommandCheck{
		Name:     "echo version",
		Command:  "echo v1.2.3",
		Expected: ">= 1.0.0",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestCommandCheck_VersionConstraintFails(t *testing.T) {
	check := &CommandCheck{
		Name:     "echo version",
		Command:  "echo v0.9.0",
		Expected: ">= 1.0.0",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestCommandCheck_VersionRange(t *testing.T) {
	check := &CommandCheck{
		Name:     "version in range",
		Command:  "echo v1.21.5",
		Expected: ">= 1.21, < 2.0",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestCommandCheck_PessimisticConstraint(t *testing.T) {
	check := &CommandCheck{
		Name:     "pessimistic",
		Command:  "echo v1.21.5",
		Expected: "~> 1.21",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v1.2.3", "1.2.3"},
		{"V1.2.3", "1.2.3"},
		{"1.2.3", "1.2.3"},
		{"  v22.3.0  ", "22.3.0"},
		{"go version go1.26.2 darwin/arm64", "1.26.2"},
		{"node v22.3.0", "22.3.0"},
		{"Python 3.12.1", "3.12.1"},
		{"rustc 1.77.0 (aedd173a2 2024-03-17)", "1.77.0"},
		{"openjdk version \"21.0.1\"", "21.0.1"},
	}

	for _, tt := range tests {
		got := extractVersion(tt.input)
		if got != tt.want {
			t.Errorf("extractVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMatchesConstraint(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		want       bool
	}{
		{"v22.3.0", ">= 22.0.0", true},
		{"v22.3.0", ">= 23.0.0", false},
		{"v1.0.0", ">= 1.0.0", true},
		{"v0.9.0", ">= 1.0.0", false},
		{"go version go1.21.0 darwin/arm64", ">= 1.21, < 2.0", true},
		{"node v18.0.0", ">= 18.0.0", true},
		{"v1.5.0", "~> 1.4", true},
		{"v2.0.0", "~> 1.4", false},
		{"v1.0.0", "!= 1.0.1", true},
	}

	for _, tt := range tests {
		got := matchesConstraint(tt.version, tt.constraint)
		if got != tt.want {
			t.Errorf("matchesConstraint(%q, %q) = %v, want %v", tt.version, tt.constraint, got, tt.want)
		}
	}
}
