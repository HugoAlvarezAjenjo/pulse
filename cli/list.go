package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
	"github.com/HugoAlvarezAjenjo/pulse/internal/styles"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured checks without running them",
	Long:  "Displays all checks defined in the pulse configuration file.",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	output := flagOutput
	if flagPlain {
		output = "plain"
	}

	switch output {
	case "json":
		return listJSON(cfg)
	case "plain":
		return listPlain(cfg)
	default:
		return listPretty(cfg)
	}
}

func listPretty(cfg *config.Config) error {
	// Calculate column widths
	maxName := 4 // "NAME"
	maxType := 4 // "TYPE"
	maxGroup := 6 // "GROUPS"
	for _, c := range cfg.Checks {
		if len(c.Name) > maxName {
			maxName = len(c.Name)
		}
		if len(c.Type) > maxType {
			maxType = len(c.Type)
		}
		g := formatGroups(c.Groups)
		if len(g) > maxGroup {
			maxGroup = len(g)
		}
	}

	// Header
	fmt.Println()
	header := fmt.Sprintf("  %-*s  %-*s  %-*s  %s", maxName, "NAME", maxType, "TYPE", maxGroup, "GROUPS", "TIMEOUT")
	fmt.Println(styles.Title.Render(header))
	fmt.Println()

	// Rows
	for _, c := range cfg.Checks {
		timeout := c.Timeout
		if timeout == "" {
			timeout = "default"
		}

		name := styles.CheckName.Render(fmt.Sprintf("%-*s", maxName, c.Name))
		typ := styles.Message.Render(fmt.Sprintf("%-*s", maxType, c.Type))
		grp := styles.Hint.Render(fmt.Sprintf("%-*s", maxGroup, formatGroups(c.Groups)))
		to := styles.Hint.Render(timeout)

		fmt.Printf("  %s  %s  %s  %s\n", name, typ, grp, to)
	}

	// Footer
	fmt.Println()
	fmt.Printf("  %s\n", styles.Message.Render(fmt.Sprintf("%d checks defined", len(cfg.Checks))))
	fmt.Println()

	return nil
}

func listPlain(cfg *config.Config) error {
	// Calculate column widths
	maxName := 4
	maxType := 4
	maxGroup := 6
	for _, c := range cfg.Checks {
		if len(c.Name) > maxName {
			maxName = len(c.Name)
		}
		if len(c.Type) > maxType {
			maxType = len(c.Type)
		}
		g := formatGroups(c.Groups)
		if len(g) > maxGroup {
			maxGroup = len(g)
		}
	}

	// Header
	fmt.Printf("%-*s  %-*s  %-*s  %s\n", maxName, "NAME", maxType, "TYPE", maxGroup, "GROUPS", "TIMEOUT")

	// Rows
	for _, c := range cfg.Checks {
		timeout := c.Timeout
		if timeout == "" {
			timeout = "default"
		}
		fmt.Printf("%-*s  %-*s  %-*s  %s\n", maxName, c.Name, maxType, c.Type, maxGroup, formatGroups(c.Groups), timeout)
	}

	fmt.Printf("\n%d checks defined\n", len(cfg.Checks))
	return nil
}

func listJSON(cfg *config.Config) error {
	type checkInfo struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Groups  []string `json:"groups,omitempty"`
		Timeout string   `json:"timeout"`
	}

	type listOutput struct {
		Checks []checkInfo `json:"checks"`
		Total  int         `json:"total"`
	}

	out := listOutput{
		Checks: make([]checkInfo, 0, len(cfg.Checks)),
		Total:  len(cfg.Checks),
	}

	for _, c := range cfg.Checks {
		timeout := c.Timeout
		if timeout == "" {
			timeout = "30s"
		}
		out.Checks = append(out.Checks, checkInfo{
			Name:    c.Name,
			Type:    c.Type,
			Groups:  c.Groups,
			Timeout: timeout,
		})
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// formatGroups returns a display string for a check's groups.
func formatGroups(groups []string) string {
	if len(groups) == 0 {
		return "*"
	}
	return strings.Join(groups, ",")
}
