package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// detection represents a detected technology and its associated checks.
type detection struct {
	Name   string
	Checks []checkEntry
}

// checkEntry represents a single check to include in the generated config.
type checkEntry struct {
	Name     string
	Type     string
	Command  string
	Expected string
	Path     string
	Host     string
	Port     int
	Variable string
	URL      string
	Status   int
	FixRun   string
}

// detectionRule defines what to look for and what checks to generate.
type detectionRule struct {
	// Files to look for (any match triggers the rule)
	Files []string
	// Generated detection result
	Detection detection
}

// rules defines all detection rules in priority order.
var rules = []detectionRule{
	{
		Files: []string{"go.mod"},
		Detection: detection{
			Name: "Go",
			Checks: []checkEntry{
				{Name: "Go", Type: "command", Command: "go version", Expected: ">=1.21"},
				{Name: "Go module", Type: "file", Path: "go.mod"},
			},
		},
	},
	{
		Files: []string{"package.json"},
		Detection: detection{
			Name: "Node.js",
			Checks: []checkEntry{
				{Name: "Node.js", Type: "command", Command: "node --version", Expected: ">=18"},
				{Name: "package.json", Type: "file", Path: "package.json"},
			},
		},
	},
	{
		Files: []string{"requirements.txt", "pyproject.toml", "setup.py"},
		Detection: detection{
			Name: "Python",
			Checks: []checkEntry{
				{Name: "Python", Type: "command", Command: "python3 --version", Expected: ">=3.10"},
				{Name: "pip", Type: "command", Command: "pip3 --version"},
			},
		},
	},
	{
		Files: []string{"Cargo.toml"},
		Detection: detection{
			Name: "Rust",
			Checks: []checkEntry{
				{Name: "Rust", Type: "command", Command: "rustc --version", Expected: ">=1.70"},
				{Name: "Cargo.toml", Type: "file", Path: "Cargo.toml"},
			},
		},
	},
	{
		Files: []string{"pom.xml", "build.gradle", "build.gradle.kts"},
		Detection: detection{
			Name: "Java",
			Checks: []checkEntry{
				{Name: "Java", Type: "command", Command: "java --version", Expected: ">=17"},
			},
		},
	},
	{
		Files: []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"},
		Detection: detection{
			Name: "Docker Compose",
			Checks: []checkEntry{
				{Name: "Docker", Type: "command", Command: "docker --version"},
			},
		},
	},
	{
		Files: []string{".env.example"},
		Detection: detection{
			Name: "Environment",
			Checks: []checkEntry{
				{Name: "Environment file", Type: "file", Path: ".env", FixRun: "cp .env.example .env"},
			},
		},
	},
	{
		Files: []string{"Makefile"},
		Detection: detection{
			Name: "Make",
			Checks: []checkEntry{
				{Name: "Make", Type: "command", Command: "make --version"},
			},
		},
	},
}

// Detect scans the given directory and generates a pulse config based on what it finds.
// Returns the YAML content and the list of detected technologies.
func Detect(dir string) (string, []string) {
	var allChecks []checkEntry
	var detected []string

	for _, rule := range rules {
		if matchesAnyFile(dir, rule.Files) {
			detected = append(detected, rule.Detection.Name)
			allChecks = append(allChecks, rule.Detection.Checks...)
		}
	}

	if len(allChecks) == 0 {
		return EmptyTemplate, nil
	}

	return renderChecks(allChecks), detected
}

// matchesAnyFile checks if any of the given files exist in the directory.
func matchesAnyFile(dir string, files []string) bool {
	for _, f := range files {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

// renderChecks generates YAML from a list of check entries.
func renderChecks(checks []checkEntry) string {
	var b strings.Builder
	b.WriteString("checks:\n")

	for _, c := range checks {
		b.WriteString(fmt.Sprintf("  - name: %s\n", c.Name))
		b.WriteString(fmt.Sprintf("    type: %s\n", c.Type))

		switch c.Type {
		case "command":
			b.WriteString(fmt.Sprintf("    command: %s\n", c.Command))
			if c.Expected != "" {
				b.WriteString(fmt.Sprintf("    expected: \"%s\"\n", c.Expected))
			}
		case "file":
			b.WriteString(fmt.Sprintf("    path: %s\n", c.Path))
		case "port":
			if c.Host != "" {
				b.WriteString(fmt.Sprintf("    host: %s\n", c.Host))
			}
			b.WriteString(fmt.Sprintf("    port: %d\n", c.Port))
		case "env":
			b.WriteString(fmt.Sprintf("    variable: %s\n", c.Variable))
		case "http":
			b.WriteString(fmt.Sprintf("    url: %s\n", c.URL))
			if c.Status != 0 {
				b.WriteString(fmt.Sprintf("    status: %d\n", c.Status))
			}
		}

		if c.FixRun != "" {
			b.WriteString("    fix:\n")
			b.WriteString(fmt.Sprintf("      run: \"%s\"\n", c.FixRun))
		}

		b.WriteString("\n")
	}

	return b.String()
}
