package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	flagPlain   bool
	flagConfig  string
	flagFix     bool
	flagOutput  string
	flagQuiet   bool
	flagTimeout string
	flagGroups  []string
)

// rootCmd is the base command for pulse.
var rootCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Validate your development environment",
	Long: `Pulse validates whether your development environment is correctly configured.

It checks for required tools, files, services, and environment variables,
then reports what's missing and optionally suggests fixes.

Run 'pulse' with no arguments to execute all checks.
Run 'pulse init' to generate a config for your project.`,
	// Default behavior: run check
	RunE: runCheck,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&flagPlain, "plain", false, "disable colors and styling (shorthand for --output plain)")
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "path to config file (default: auto-discover)")
	rootCmd.PersistentFlags().BoolVar(&flagFix, "fix", false, "prompt to run suggested fixes for failed checks")
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "pretty", "output format: pretty, plain, json, github")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "only show failures and errors")
	rootCmd.PersistentFlags().StringVarP(&flagTimeout, "timeout", "t", "30s", "timeout per check (e.g., 10s, 1m)")
	rootCmd.PersistentFlags().StringSliceVarP(&flagGroups, "group", "g", nil, "run only checks in specified groups (can be repeated)")

	// Add subcommands
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(ciCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(completionCmd)
}

// SetVersion configures the version information for the CLI.
func SetVersion(version, commit, date string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{printf "pulse %s\n  commit: %s\n  built:  %s\n" .Version "` + commit + `" "` + date + `"}}`)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
