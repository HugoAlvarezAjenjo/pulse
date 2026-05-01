package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const defaultConfigTemplate = `checks:
  - name: Go
    type: command
    command: go version
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a starter pulse config",
	Long:  "Generates a minimal .pulse.yaml configuration file in the current directory.",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	path := filepath.Join(dir, ".pulse.yaml")

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "error: %s already exists\n", path)
		os.Exit(1)
	}

	if err := os.WriteFile(path, []byte(defaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Printf("Created %s\n", path)
	return nil
}
