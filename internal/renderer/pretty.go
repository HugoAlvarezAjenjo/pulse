package renderer

import (
	"fmt"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
	"github.com/HugoAlvarezAjenjo/pulse/internal/styles"
)

// PrettyRenderer outputs styled terminal output using Lip Gloss.
type PrettyRenderer struct{}

// NewPretty creates a new PrettyRenderer.
func NewPretty() *PrettyRenderer {
	return &PrettyRenderer{}
}

func (p *PrettyRenderer) Header() {
	fmt.Println()
	fmt.Println(styles.Title.Render("Pulse environment check"))
	fmt.Println()
}

func (p *PrettyRenderer) Success(r result.Result) {
	icon := styles.SuccessIcon.String()
	name := styles.CheckName.Render(r.Name)
	msg := ""
	if r.Message != "" {
		msg = " " + styles.Message.Render(r.Message)
	}
	fmt.Printf("%s %s%s\n", icon, name, msg)
}

func (p *PrettyRenderer) Failure(r result.Result) {
	icon := styles.FailureIcon.String()
	name := styles.CheckName.Render(r.Name)
	fmt.Printf("%s %s\n", icon, name)

	if r.Message != "" {
		fmt.Printf("  %s\n", styles.Message.Render(r.Message))
	}
	if r.Hint != "" {
		fmt.Printf("  %s\n", styles.Hint.Render(r.Hint))
	}
	if r.Fix != nil {
		fmt.Println()
		fmt.Printf("  %s\n", styles.FixLabel.Render("Suggested fix:"))
		fmt.Printf("  %s\n", styles.FixCommand.Render(r.Fix.Run))
	}
	fmt.Println()
}

func (p *PrettyRenderer) Error(r result.Result) {
	icon := styles.ErrorIcon.String()
	name := styles.CheckName.Render(r.Name)
	fmt.Printf("%s %s\n", icon, name)

	if r.Message != "" {
		fmt.Printf("  %s\n", styles.Message.Render(r.Message))
	}
	fmt.Println()
}

func (p *PrettyRenderer) Summary(results []result.Result) {
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
	fmt.Print("Summary: ")

	parts := []string{}
	if passed > 0 {
		parts = append(parts, styles.SummaryPassed.Render(fmt.Sprintf("%d passed", passed)))
	}
	if failed > 0 {
		parts = append(parts, styles.SummaryFailed.Render(fmt.Sprintf("%d failed", failed)))
	}
	if errors > 0 {
		parts = append(parts, styles.SummaryFailed.Render(fmt.Sprintf("%d errors", errors)))
	}

	dot := " " + styles.SummaryDot.String() + " "
	for i, part := range parts {
		if i > 0 {
			fmt.Print(dot)
		}
		fmt.Print(part)
	}
	fmt.Println()
	fmt.Println()
}
