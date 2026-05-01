package runner

import (
	"context"
	"testing"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/checks"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// mockCheck is a simple check implementation for testing.
type mockCheck struct {
	name   string
	status result.Status
	delay  time.Duration
}

func (m *mockCheck) Run(ctx context.Context) result.Result {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return result.Result{
				Name:   m.name,
				Status: result.Error,
			}
		}
	}
	return result.Result{
		Name:   m.name,
		Status: m.status,
	}
}

func TestRunner_PreservesOrder(t *testing.T) {
	checkList := []checks.Check{
		&mockCheck{name: "first", status: result.Success, delay: 50 * time.Millisecond},
		&mockCheck{name: "second", status: result.Failure, delay: 10 * time.Millisecond},
		&mockCheck{name: "third", status: result.Success, delay: 30 * time.Millisecond},
	}

	r := New()
	results := r.Run(context.Background(), checkList)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	expected := []string{"first", "second", "third"}
	for i, name := range expected {
		if results[i].Name != name {
			t.Errorf("result[%d]: expected name %q, got %q", i, name, results[i].Name)
		}
	}
}

func TestRunner_ConcurrentExecution(t *testing.T) {
	// All checks have 50ms delay; if sequential would take 150ms+
	checkList := []checks.Check{
		&mockCheck{name: "a", status: result.Success, delay: 50 * time.Millisecond},
		&mockCheck{name: "b", status: result.Success, delay: 50 * time.Millisecond},
		&mockCheck{name: "c", status: result.Success, delay: 50 * time.Millisecond},
	}

	r := New()
	start := time.Now()
	results := r.Run(context.Background(), checkList)
	elapsed := time.Since(start)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// Should complete in ~50ms, not ~150ms
	if elapsed > 120*time.Millisecond {
		t.Errorf("expected concurrent execution (<120ms), took %s", elapsed)
	}
}

func TestRunner_Timeout(t *testing.T) {
	checkList := []checks.Check{
		&mockCheck{name: "slow", status: result.Success, delay: 5 * time.Second},
	}

	r := &Runner{Timeout: 100 * time.Millisecond}
	results := r.Run(context.Background(), checkList)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != result.Error {
		t.Errorf("expected error on timeout, got %s", results[0].Status)
	}
}

func TestRunner_EmptyChecks(t *testing.T) {
	r := New()
	results := r.Run(context.Background(), nil)

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
