package renderer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// JSONOutput represents the full JSON output structure.
type JSONOutput struct {
	Results []JSONResult `json:"results"`
	Summary JSONSummary  `json:"summary"`
}

// JSONResult represents a single check result in JSON.
type JSONResult struct {
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	Message    string  `json:"message,omitempty"`
	Hint       string  `json:"hint,omitempty"`
	Fix        string  `json:"fix,omitempty"`
	DurationMs int64   `json:"duration_ms"`
}

// JSONSummary represents the summary section in JSON output.
type JSONSummary struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
	Errors int `json:"errors"`
}

// JSONRenderer outputs results as structured JSON.
type JSONRenderer struct {
	Quiet   bool
	results []result.Result
}

// NewJSON creates a new JSONRenderer.
func NewJSON(quiet bool) *JSONRenderer {
	return &JSONRenderer{Quiet: quiet}
}

func (j *JSONRenderer) Header() {
	// No header in JSON mode
}

func (j *JSONRenderer) Success(r result.Result) {
	j.results = append(j.results, r)
}

func (j *JSONRenderer) Failure(r result.Result) {
	j.results = append(j.results, r)
}

func (j *JSONRenderer) Error(r result.Result) {
	j.results = append(j.results, r)
}

func (j *JSONRenderer) Summary(results []result.Result) {
	output := JSONOutput{
		Results: make([]JSONResult, 0, len(results)),
		Summary: JSONSummary{Total: len(results)},
	}

	for _, r := range results {
		// In quiet mode, skip successes
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

		switch r.Status {
		case result.Success:
			output.Summary.Passed++
		case result.Failure:
			output.Summary.Failed++
		case result.Error:
			output.Summary.Errors++
		}
	}

	// Count all for summary regardless of quiet
	output.Summary.Passed = 0
	output.Summary.Failed = 0
	output.Summary.Errors = 0
	for _, r := range results {
		switch r.Status {
		case result.Success:
			output.Summary.Passed++
		case result.Failure:
			output.Summary.Failed++
		case result.Error:
			output.Summary.Errors++
		}
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal JSON: %s\n", err)
		return
	}

	fmt.Println(string(data))
}
