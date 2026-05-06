package renderer

import (
	"fmt"
	"time"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// PlainRenderer outputs unformatted text suitable for CI and piping.
type PlainRenderer struct {
	Quiet bool
}

// NewPlain creates a new PlainRenderer.
func NewPlain(quiet bool) *PlainRenderer {
	return &PlainRenderer{Quiet: quiet}
}

func (p *PlainRenderer) Header() {
	if p.Quiet {
		return
	}
	fmt.Println("Pulse environment check")
	fmt.Println()
}

func (p *PlainRenderer) Success(r result.Result) {
	if p.Quiet {
		return
	}
	msg := ""
	if r.Message != "" {
		msg = " " + r.Message
	}
	fmt.Printf("[PASS] %s%s\n", r.Name, msg)
}

func (p *PlainRenderer) Failure(r result.Result) {
	fmt.Printf("[FAIL] %s\n", r.Name)
	if r.Message != "" {
		fmt.Printf("       %s\n", r.Message)
	}
	if r.Hint != "" {
		fmt.Printf("       %s\n", r.Hint)
	}
	if r.Fix != nil {
		fmt.Printf("       fix: %s\n", r.Fix.Run)
	}
}

func (p *PlainRenderer) Error(r result.Result) {
	fmt.Printf("[ERR]  %s\n", r.Name)
	if r.Message != "" {
		fmt.Printf("       %s\n", r.Message)
	}
}

func (p *PlainRenderer) Summary(results []result.Result, duration time.Duration) {
	passed := 0
	failed := 0
	errors := 0

	for _, r := range results {
		switch r.Status {
		case result.Success:
			passed++
		case result.Failure:
			failed++
		case result.Error:
			errors++
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d passed, %d failed, %d errors in %s\n", passed, failed, errors, formatDuration(duration))
}
