package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/checks"
	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
	"github.com/HugoAlvarezAjenjo/pulse/internal/fix"
	"github.com/HugoAlvarezAjenjo/pulse/internal/renderer"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
	"github.com/HugoAlvarezAjenjo/pulse/internal/runner"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run environment checks",
	Long:  "Validates the development environment against the pulse configuration.",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
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

	// Select renderer
	rnd := selectRenderer()

	// Run checks
	r := runner.New()
	if timeout, err := time.ParseDuration(flagTimeout); err == nil {
		r.Timeout = timeout
	}
	results := r.Run(ctx, checkList)

	// Render results
	renderer.Render(rnd, results)

	// Handle fixes if enabled (not in json/github mode)
	if flagFix && flagOutput != "json" && flagOutput != "github" {
		handleFixes(ctx, results)
	}

	// Exit with appropriate code
	os.Exit(exitCode(results))
	return nil
}

func loadConfig() (*config.Config, error) {
	if flagConfig != "" {
		return config.Load(flagConfig)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}

	cfg, _, err := config.DiscoverAndLoad(dir)
	return cfg, err
}

func buildChecks(cfg *config.Config) ([]checks.Check, error) {
	checkList := make([]checks.Check, 0, len(cfg.Checks))

	for _, checkCfg := range cfg.Checks {
		c, err := checks.FromConfig(checkCfg)
		if err != nil {
			return nil, err
		}
		checkList = append(checkList, c)
	}

	return checkList, nil
}

func selectRenderer() renderer.Renderer {
	// --plain is shorthand for --output plain
	output := flagOutput
	if flagPlain {
		output = "plain"
	}

	switch output {
	case "json":
		return renderer.NewJSON(flagQuiet)
	case "github":
		return renderer.NewGitHub(flagQuiet)
	case "plain":
		return renderer.NewPlain(flagQuiet)
	default:
		return renderer.NewPretty(flagQuiet)
	}
}

func handleFixes(ctx context.Context, results []result.Result) {
	executor := fix.New()

	for _, r := range results {
		if r.Status == result.Failure && r.Fix != nil {
			executor.PromptAndRun(ctx, r)
		}
	}
}

func exitCode(results []result.Result) int {
	for _, r := range results {
		switch r.Status {
		case result.Failure:
			return 1
		case result.Error:
			return 2
		}
	}
	return 0
}
