package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/checks"
	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
	"github.com/HugoAlvarezAjenjo/pulse/internal/styles"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate config syntax without running checks",
	Long:  "Parses and validates the pulse configuration file, reporting any errors without executing checks.",
	RunE:  runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	// Discover config
	configPath := flagConfig
	if configPath == "" {
		path, discErr := config.Discover(dir)
		if discErr != nil {
			fmt.Fprintf(os.Stderr, "%s No config file found\n", styles.FailureIcon.String())
			os.Exit(2)
		}
		configPath = path
	}

	// Load and parse
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("%s %s\n", styles.FailureIcon.String(), styles.Message.Render(fmt.Sprintf("config error: %s", err)))
		os.Exit(2)
	}

	// Validate each check can be constructed
	var validationErrors []string
	for _, checkCfg := range cfg.Checks {
		if _, err := checks.FromConfig(checkCfg); err != nil {
			validationErrors = append(validationErrors, err.Error())
		}

		// Validate timeout format if specified
		if checkCfg.Timeout != "" {
			if _, err := time.ParseDuration(checkCfg.Timeout); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("check %q: invalid timeout %q", checkCfg.Name, checkCfg.Timeout))
			}
		}

		// Validate groups are non-empty strings
		for _, g := range checkCfg.Groups {
			if g == "" {
				validationErrors = append(validationErrors, fmt.Sprintf("check %q: empty group name", checkCfg.Name))
			}
		}
	}

	if len(validationErrors) > 0 {
		fmt.Printf("%s Config has %d error(s):\n\n", styles.FailureIcon.String(), len(validationErrors))
		for _, e := range validationErrors {
			fmt.Printf("  • %s\n", styles.Message.Render(e))
		}
		fmt.Println()
		os.Exit(1)
	}

	// Check for local overrides
	localInfo := ""
	if config.DiscoverLocal(dir) != "" {
		localInfo = " (+ .pulse.local.yaml)"
	}

	fmt.Printf("%s %s is valid%s — %d checks defined\n",
		styles.SuccessIcon.String(),
		styles.CheckName.Render(configPath),
		localInfo,
		len(cfg.Checks),
	)

	return nil
}
