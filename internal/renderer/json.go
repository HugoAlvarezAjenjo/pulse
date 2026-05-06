package renderer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// JSONOutput represents the full JSON output structure.
type JSONOutput struct {
	Results []JSONResult `json:"results"`
	Summary JSONSummary  `json:"summary"`
}

// JSONResult represents a single check result in JSON.
type JSONResult struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
	Hint       string `json:"hint,omitempty"`
	Fix        string `json:"fix,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

// JSONSummary represents the summary section in JSON output.
type JSONSummary struct {
	Total      int   `json:"total"`
	Passed     int   `json:"passed"`
	Failed     int   `json:"failed"`
	Errors     int   `json:"errors"`
	DurationMs int64 `json:"duration_ms"`
}

// JSONRenderer outputs results as structured JSON.
type JSONRenderer struct {
	Quiet bool
}

// NewJSON creates a new JSONRenderer.
func NewJSON(quiet bool) *JSONRenderer {
	return &JSONRenderer{Quiet: quiet}
}

func (j *JSONRenderer) Header() {}

func (j *JSONRenderer) Success(_ result.Result) {}

func (j *JSONRenderer) Failure(_ result.Result) {}

func (j *JSONRenderer) Error(_ result.Result) {}

func (j *JSONRenderer) Summary(results []result.Result, duration time.Duration) {
	output := JSONOutput{
		Results: make([]JSONResult, 0, len(results)),
		Summary: JSONSummary{Total: len(results), DurationMs: duration.Milliseconds()},
	}

	for _, r := range results {
		// Count for summary
		switch r.Status {
		case result.Success:
			output.Summary.Passed++
		case result.Failure:
			output.Summary.Failed++
		case result.Error:
			output.Summary.Errors++
		}

		// In quiet mode, skip successes from results list
		if j.Quiet && r.Status == result.Success {
			continue
		}

		jr := JSONResult{
			Name:       r.Name,
			Status:     r.Status.String(),
			Message:    r.Message,
			Hint:       r.Hint,
			DurationMs: r.Duration.Milliseconds(),
		}
		if r.Fix != nil {
			jr.Fix = r.Fix.Run
		}

		output.Results = append(output.Results, jr)
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal JSON: %s\n", err)
		return
	}

	fmt.Println(string(data))
}
