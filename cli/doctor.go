package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
	"github.com/HugoAlvarezAjenjo/pulse/internal/styles"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose pulse installation and configuration",
	Long:  "Shows diagnostic information about the pulse installation, environment, and configuration.",
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	label := styles.CheckName
	value := styles.Message

	fmt.Println()
	fmt.Println(styles.Title.Render("Pulse Doctor"))
	fmt.Println()

	// Version info
	fmt.Printf("  %-10s %s\n", label.Render("Version:"), value.Render(rootCmd.Version))
	fmt.Printf("  %-10s %s\n", label.Render("OS:"), value.Render(runtime.GOOS+"/"+runtime.GOARCH))
	fmt.Printf("  %-10s %s\n", label.Render("Go:"), value.Render(runtime.Version()))

	// Working directory
	dir, err := os.Getwd()
	if err != nil {
		dir = "unknown"
	}
	fmt.Printf("  %-10s %s\n", label.Render("Workdir:"), value.Render(dir))

	fmt.Println()

	// Config discovery
	fmt.Println(styles.Title.Render("  Configuration"))
	fmt.Println()

	configPath, discoverErr := config.Discover(dir)
	if discoverErr != nil {
		fmt.Printf("  %s %s\n", styles.FailureIcon.String(), value.Render("no config file found"))
		fmt.Printf("    %s\n", styles.Hint.Render("run 'pulse init' to create one"))
	} else {
		fmt.Printf("  %s %s\n", styles.SuccessIcon.String(), value.Render(configPath))

		// Try loading config to check validity
		cfg, loadErr := config.Load(configPath)
		if loadErr != nil {
			fmt.Printf("  %s %s\n", styles.FailureIcon.String(), styles.Message.Render(fmt.Sprintf("config error: %s", loadErr)))
		} else {
			fmt.Printf("  %s %s\n", styles.SuccessIcon.String(), value.Render(fmt.Sprintf("%d checks defined", len(cfg.Checks))))

			// Show check types summary
			types := countCheckTypes(cfg)
			if len(types) > 0 {
				fmt.Printf("    %s\n", styles.Hint.Render(formatTypes(types)))
			}
		}
	}

	fmt.Println()
	return nil
}

func countCheckTypes(cfg *config.Config) map[string]int {
	types := make(map[string]int)
	for _, c := range cfg.Checks {
		types[c.Type]++
	}
	return types
}

func formatTypes(types map[string]int) string {
	result := ""
	for t, count := range types {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%s: %d", t, count)
	}
	return result
}
