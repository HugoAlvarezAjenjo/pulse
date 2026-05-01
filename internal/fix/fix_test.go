package fix

import (
	"context"
	"testing"
	"time"
)

func TestExecutor_RunSuccess(t *testing.T) {
	e := &Executor{Timeout: 5 * time.Second}
	err := e.Run(context.Background(), "echo hello")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
}

func TestExecutor_RunFailure(t *testing.T) {
	e := &Executor{Timeout: 5 * time.Second}
	err := e.Run(context.Background(), "exit 1")
	if err == nil {
		t.Error("expected error for failed command")
	}
}

func TestExecutor_RunTimeout(t *testing.T) {
	e := &Executor{Timeout: 100 * time.Millisecond}
	err := e.Run(context.Background(), "sleep 10")
	if err == nil {
		t.Error("expected error for timed out command")
	}
}
