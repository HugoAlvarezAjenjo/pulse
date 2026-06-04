package renderer

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPlainRenderer_Output(t *testing.T) {
	results := []result.Result{
		{Name: "Go", Status: result.Success, Message: "1.26.2", Duration: 10 * time.Millisecond},
		{Name: "DB", Status: result.Failure, Message: "connection refused", Duration: 5 * time.Millisecond},
	}

	output := captureStdout(t, func() {
		rnd := NewPlain(false)
		Render(rnd, results, 15*time.Millisecond)
	})

	if !strings.Contains(output, "PASS") {
		t.Error("plain output should contain PASS for success")
	}
	if !strings.Contains(output, "FAIL") {
		t.Error("plain output should contain FAIL for failure")
	}
	if !strings.Contains(output, "Go") {
		t.Error("plain output should contain check name")
	}
}

func TestPlainRenderer_Quiet(t *testing.T) {
	results := []result.Result{
		{Name: "Go", Status: result.Success, Message: "1.26.2", Duration: 10 * time.Millisecond},
		{Name: "DB", Status: result.Failure, Message: "refused", Duration: 5 * time.Millisecond},
	}

	output := captureStdout(t, func() {
		rnd := NewPlain(true)
		Render(rnd, results, 15*time.Millisecond)
	})

	if strings.Contains(output, "PASS") {
		t.Error("quiet mode should not show passing checks")
	}
	if !strings.Contains(output, "FAIL") {
		t.Error("quiet mode should still show failures")
	}
}

func TestJSONRenderer_ValidJSON(t *testing.T) {
	results := []result.Result{
		{Name: "Go", Status: result.Success, Message: "1.26.2", Duration: 10 * time.Millisecond},
		{Name: "DB", Status: result.Failure, Message: "refused", Duration: 5 * time.Millisecond,
			Fix: &result.Fix{Run: "docker compose up db -d"}},
	}

	output := captureStdout(t, func() {
		rnd := NewJSON(false)
		Render(rnd, results, 15*time.Millisecond)
	})

	var parsed JSONOutput
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("JSON output is not valid: %v\nOutput: %s", err, output)
	}

	if parsed.Summary.Total != 2 {
		t.Errorf("expected total=2, got %d", parsed.Summary.Total)
	}
	if parsed.Summary.Passed != 1 {
		t.Errorf("expected passed=1, got %d", parsed.Summary.Passed)
	}
	if parsed.Summary.Failed != 1 {
		t.Errorf("expected failed=1, got %d", parsed.Summary.Failed)
	}
}

func TestJSONRenderer_Quiet_HidesSuccess(t *testing.T) {
	results := []result.Result{
		{Name: "Go", Status: result.Success, Duration: 10 * time.Millisecond},
		{Name: "DB", Status: result.Failure, Duration: 5 * time.Millisecond},
	}

	output := captureStdout(t, func() {
		rnd := NewJSON(true)
		Render(rnd, results, 15*time.Millisecond)
	})

	var parsed JSONOutput
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("JSON output is not valid: %v", err)
	}

	if len(parsed.Results) != 1 {
		t.Errorf("quiet mode: expected 1 result (failure only), got %d", len(parsed.Results))
	}
	if parsed.Summary.Total != 2 {
		t.Errorf("summary should still count all, got total=%d", parsed.Summary.Total)
	}
}

func TestGitHubRenderer_Annotations(t *testing.T) {
	results := []result.Result{
		{Name: "Go", Status: result.Success, Message: "1.26.2", Duration: 10 * time.Millisecond},
		{Name: "DB", Status: result.Failure, Message: "refused", Duration: 5 * time.Millisecond},
	}

	output := captureStdout(t, func() {
		rnd := NewGitHub(false)
		Render(rnd, results, 15*time.Millisecond)
	})

	if !strings.Contains(output, "::notice title=Go::") {
		t.Error("expected GitHub notice annotation for success")
	}
	if !strings.Contains(output, "::error title=DB::") {
		t.Error("expected GitHub error annotation for failure")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{500 * time.Microsecond, "500µs"},
		{50 * time.Millisecond, "50ms"},
		{1500 * time.Millisecond, "1.5s"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.d)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
