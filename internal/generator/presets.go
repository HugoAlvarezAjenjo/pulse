package generator

// Preset templates for common stacks.
var presets = map[string]string{
	"go": `checks:
  - name: Go
    type: command
    command: go version
    expected: ">= 1.21"

  - name: Go module
    type: file
    path: go.mod

  - name: Main entry point
    type: file
    path: main.go
`,

	"node": `checks:
  - name: Node.js
    type: command
    command: node --version
    expected: ">= 18.0"

  - name: npm
    type: command
    command: npm --version

  - name: package.json
    type: file
    path: package.json

  - name: Environment file
    type: file
    path: .env
    fix:
      run: "cp .env.example .env"
`,

	"python": `checks:
  - name: Python
    type: command
    command: python3 --version
    expected: ">= 3.10"

  - name: pip
    type: command
    command: pip3 --version

  - name: Requirements
    type: file
    path: requirements.txt

  - name: Virtual environment
    type: file
    path: .venv
    fix:
      run: "python3 -m venv .venv"
`,

	"java": `checks:
  - name: Java
    type: command
    command: java --version
    expected: ">= 17.0"

  - name: Maven/Gradle
    type: file
    path: pom.xml

  - name: Source directory
    type: file
    path: src/main/java
`,

	"rust": `checks:
  - name: Rust
    type: command
    command: rustc --version
    expected: ">= 1.70"

  - name: Cargo
    type: command
    command: cargo --version

  - name: Cargo.toml
    type: file
    path: Cargo.toml
`,
}

// EmptyTemplate is the blank template with commented guidance.
const EmptyTemplate = `checks:
  # Add your checks here. Examples:
  #
  # - name: Node.js
  #   type: command
  #   command: node --version
  #   expected: ">=18"
  #
  # - name: Config file
  #   type: file
  #   path: .env
  #
  # - name: Database
  #   type: port
  #   host: localhost
  #   port: 5432
  #   fix:
  #     run: "docker compose up db -d"
  #
  # Check types: command, file, port, env, http, docker
  # Docs: https://github.com/HugoAlvarezAjenjo/pulse

  - name: Example
    type: command
    command: echo "pulse is ready"
`

// GetPreset returns the template for a given preset name.
// Returns empty string if not found.
func GetPreset(name string) (string, bool) {
	t, ok := presets[name]
	return t, ok
}

// AvailablePresets returns the list of available preset names.
func AvailablePresets() []string {
	names := make([]string, 0, len(presets))
	for name := range presets {
		names = append(names, name)
	}
	return names
}
