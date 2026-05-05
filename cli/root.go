package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	flagPlain  bool
	flagConfig string
	flagFix    bool
	flagOutput string
	flagQuiet  bool
)

// rootCmd is the base command for pulse.
var rootCmd = &cobra.Command{
	Use:   "pulse",
	Short: "Takes the pulse of your development environment",
	Long:  "Pulse validates whether your development environment is correctly configured for a project.",
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

	// Add subcommands
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(doctorCmd)
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
