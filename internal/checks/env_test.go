package checks

import (
	"context"
	"os"
	"testing"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestEnvCheck_Exists(t *testing.T) {
	t.Setenv("PULSE_TEST_VAR", "hello")

	check := &EnvCheck{
		Name:     "test env",
		Variable: "PULSE_TEST_VAR",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestEnvCheck_NotSet(t *testing.T) {
	os.Unsetenv("PULSE_TEST_MISSING_VAR") //nolint:errcheck

	check := &EnvCheck{
		Name:     "missing env",
		Variable: "PULSE_TEST_MISSING_VAR",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestEnvCheck_Empty(t *testing.T) {
	t.Setenv("PULSE_TEST_EMPTY", "")

	check := &EnvCheck{
		Name:     "empty env",
		Variable: "PULSE_TEST_EMPTY",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure for empty var, got %s", r.Status)
	}
}

func TestEnvCheck_PatternMatch(t *testing.T) {
	t.Setenv("PULSE_TEST_KEY", "sk-abc123")

	check := &EnvCheck{
		Name:     "api key",
		Variable: "PULSE_TEST_KEY",
		Expected: "sk-*",
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestEnvCheck_PatternNoMatch(t *testing.T) {
	t.Setenv("PULSE_TEST_KEY", "wrong-value")

	check := &EnvCheck{
		Name:     "api key",
		Variable: "PULSE_TEST_KEY",
		Expected: "sk-*",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		value   string
		pattern string
		want    bool
	}{
		{"sk-abc123", "sk-*", true},
		{"wrong", "sk-*", false},
		{"hello-production", "*-production", true},
		{"hello-staging", "*-production", false},
		{"exact", "exact", true},
		{"different", "exact", false},
	}

	for _, tt := range tests {
		got := matchPattern(tt.value, tt.pattern)
		if got != tt.want {
			t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.value, tt.pattern, got, tt.want)
		}
	}
}
