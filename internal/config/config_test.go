package config

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	yaml := `
checks:
  - name: Node.js
    type: command
    command: node --version
    expected: ">=22"
  - name: Environment file
    type: file
    path: .env
  - name: PostgreSQL
    type: port
    host: localhost
    port: 5432
    fix:
      run: "docker compose up db -d"
`

	cfg, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(cfg.Checks) != 3 {
		t.Fatalf("expected 3 checks, got %d", len(cfg.Checks))
	}

	// Check first entry
	if cfg.Checks[0].Name != "Node.js" {
		t.Errorf("expected name 'Node.js', got %q", cfg.Checks[0].Name)
	}
	if cfg.Checks[0].Type != "command" {
		t.Errorf("expected type 'command', got %q", cfg.Checks[0].Type)
	}
	if cfg.Checks[0].Command != "node --version" {
		t.Errorf("expected command 'node --version', got %q", cfg.Checks[0].Command)
	}
	if cfg.Checks[0].Expected != ">=22" {
		t.Errorf("expected '>=22', got %q", cfg.Checks[0].Expected)
	}

	// Check fix
	if cfg.Checks[2].Fix == nil {
		t.Fatal("expected fix to be set")
	}
	if cfg.Checks[2].Fix.Run != "docker compose up db -d" {
		t.Errorf("expected fix run 'docker compose up db -d', got %q", cfg.Checks[2].Fix.Run)
	}
}

func TestParse_Empty(t *testing.T) {
	yaml := `checks: []`

	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for empty checks")
	}
}

func TestParse_MissingName(t *testing.T) {
	yaml := `
checks:
  - type: command
    command: echo hi
`

	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParse_MissingType(t *testing.T) {
	yaml := `
checks:
  - name: Test
    command: echo hi
`

	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	yaml := `invalid: [yaml: content`

	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
