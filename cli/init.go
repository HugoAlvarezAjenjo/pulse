package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/generator"
	"github.com/HugoAlvarezAjenjo/pulse/internal/styles"
)

var (
	flagPreset string
	flagEmpty  bool
	flagForce  bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a pulse config file",
	Long: `Generates a .pulse.yaml configuration file in the current directory.

Modes:
  pulse init              Auto-detect stack and generate checks
  pulse init --preset go  Use a predefined template (go, node, python, java, rust)
  pulse init --empty      Generate a blank template with guidance comments`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&flagPreset, "preset", "", "use a preset template (go, node, python, java, rust)")
	initCmd.Flags().BoolVar(&flagEmpty, "empty", false, "generate a blank template")
	initCmd.Flags().BoolVar(&flagForce, "force", false, "overwrite existing config file")
}

func runInit(cmd *cobra.Command, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	path := filepath.Join(dir, ".pulse.yaml")

	// Check if file already exists
	if !flagForce {
		if _, err := os.Stat(path); err == nil {
			fmt.Fprintf(os.Stderr, "error: %s already exists (use --force to overwrite)\n", path)
			os.Exit(1)
		}
	}

	// Determine content based on mode
	var content string
	switch {
	case flagEmpty:
		content = generator.EmptyTemplate
		writeConfig(path, content)
		fmt.Printf("%s Created %s (empty template)\n", styles.SuccessIcon.String(), path)

	case flagPreset != "":
		preset, ok := generator.GetPreset(flagPreset)
		if !ok {
			available := strings.Join(generator.AvailablePresets(), ", ")
			fmt.Fprintf(os.Stderr, "error: unknown preset %q (available: %s)\n", flagPreset, available)
			os.Exit(1)
		}
		content = preset
		writeConfig(path, content)
		fmt.Printf("%s Created %s (preset: %s)\n", styles.SuccessIcon.String(), path, flagPreset)

	default:
		// Auto-detect
		content, detected := generator.Detect(dir)
		writeConfig(path, content)

		if len(detected) > 0 {
			fmt.Printf("\n%s Detected: %s\n", styles.SuccessIcon.String(), strings.Join(detected, ", "))
			fmt.Printf("%s Created %s with %d checks\n\n", styles.SuccessIcon.String(), path, countChecksInContent(content))
		} else {
			fmt.Printf("%s No stack detected, created %s with example template\n", styles.SuccessIcon.String(), path)
			fmt.Printf("  %s\n", styles.Hint.Render("Edit the file to add your checks"))
		}
	}

	return nil
}

func writeConfig(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: writing config: %s\n", err)
		os.Exit(2)
	}
}

func countChecksInContent(content string) int {
	count := 0
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- name:") {
			count++
		}
	}
	return count
}
