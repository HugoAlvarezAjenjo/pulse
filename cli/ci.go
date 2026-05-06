package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/renderer"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
	"github.com/HugoAlvarezAjenjo/pulse/internal/runner"
)

var flagFailFast bool

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Run checks in CI mode (no colors, no prompts, fail-fast optional)",
	Long: `Runs environment checks optimized for CI/CD pipelines.

Defaults: plain output, quiet mode, no fix prompts.
Designed for non-interactive environments like GitHub Actions, GitLab CI, Jenkins.`,
	RunE: runCI,
}

func init() {
	ciCmd.Flags().BoolVar(&flagFailFast, "fail-fast", false, "stop at first failed check")
}

func runCI(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load config
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	// Build checks from config
	checkList, err := buildChecks(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	// Filter by groups if specified
	if len(flagGroups) > 0 {
		checkList = filterByGroups(cfg, checkList, flagGroups)
	}

	// CI always uses plain or the override output (json/github), never pretty
	rnd := selectCIRenderer()

	// Parse timeout
	r := runner.New()
	if timeout, err := time.ParseDuration(flagTimeout); err == nil {
		r.Timeout = timeout
	}

	// Execute
	var results []result.Result
	var duration time.Duration

	if flagFailFast {
		results, duration = runFailFast(ctx, r, checkList)
	} else {
		start := time.Now()
		results = r.RunWithTimeouts(ctx, checkList)
		duration = time.Since(start)
	}

	// Render
	renderer.Render(rnd, results, duration)

	// Exit — never prompt for fixes in CI
	os.Exit(exitCode(results))
	return nil
}

// selectCIRenderer picks the renderer for CI mode.
// Defaults to plain (quiet), but respects --output if explicitly set.
func selectCIRenderer() renderer.Renderer {
	switch flagOutput {
	case "json":
		return renderer.NewJSON(true)
	case "github":
		return renderer.NewGitHub(true)
	default:
		// CI default: plain + quiet
		return renderer.NewPlain(true)
	}
}

// runFailFast executes checks sequentially, stopping at the first failure.
func runFailFast(ctx context.Context, r *runner.Runner, checkList []runner.CheckWithTimeout) ([]result.Result, time.Duration) {
	start := time.Now()
	results := make([]result.Result, 0, len(checkList))

	for _, cwt := range checkList {
		timeout := r.Timeout
		if cwt.Timeout > 0 {
			timeout = cwt.Timeout
		}

		checkCtx, cancel := context.WithTimeout(ctx, timeout)
		res := cwt.Check.Run(checkCtx)
		cancel()

		results = append(results, res)

		if res.Status == result.Failure || res.Status == result.Error {
			break
		}
	}

	return results, time.Since(start)
}
