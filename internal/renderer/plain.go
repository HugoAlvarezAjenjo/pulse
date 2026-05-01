package renderer

import (
	"fmt"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// PlainRenderer outputs unformatted text suitable for CI and piping.
type PlainRenderer struct{}

// NewPlain creates a new PlainRenderer.
func NewPlain() *PlainRenderer {
	return &PlainRenderer{}
}

func (p *PlainRenderer) Header() {
	fmt.Println("Pulse environment check")
	fmt.Println()
}

func (p *PlainRenderer) Success(r result.Result) {
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

func (p *PlainRenderer) Summary(results []result.Result) {
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
	fmt.Printf("Summary: %d passed, %d failed, %d errors\n", passed, failed, errors)
}
