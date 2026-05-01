package renderer

import (
	"fmt"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// GitHubRenderer outputs results as GitHub Actions workflow commands.
// See: https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions
type GitHubRenderer struct {
	Quiet bool
}

// NewGitHub creates a new GitHubRenderer.
func NewGitHub(quiet bool) *GitHubRenderer {
	return &GitHubRenderer{Quiet: quiet}
}

func (g *GitHubRenderer) Header() {
	// No header in GitHub mode
}

func (g *GitHubRenderer) Success(r result.Result) {
	if g.Quiet {
		return
	}
	msg := r.Name
	if r.Message != "" {
		msg = r.Message
	}
	fmt.Printf("::notice title=%s::%s\n", r.Name, msg)
}

func (g *GitHubRenderer) Failure(r result.Result) {
	msg := r.Name
	if r.Message != "" {
		msg = r.Message
	}
	if r.Hint != "" {
		msg += " — " + r.Hint
	}
	fmt.Printf("::error title=%s::%s\n", r.Name, msg)
}

func (g *GitHubRenderer) Error(r result.Result) {
	msg := r.Name
	if r.Message != "" {
		msg = r.Message
	}
	fmt.Printf("::warning title=%s::%s\n", r.Name, msg)
}

func (g *GitHubRenderer) Summary(results []result.Result) {
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

	fmt.Printf("::notice::Pulse summary: %d passed, %d failed, %d errors\n", passed, failed, errors)
}
