package config

import (
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	yaml := `
checks:
  - name: Go
    type: command
    command: go version
  - name: Config
    type: file
    path: .env
  - name: DB
    type: port
    host: localhost
    port: 5432
`
	_, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("expected valid config, got error: %s", err)
	}
}

func TestValidate_EmptyChecks(t *testing.T) {
	yaml := `checks: []`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for empty checks")
	}
}

func TestValidate_MissingName(t *testing.T) {
	yaml := `
checks:
  - type: command
    command: go version
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestValidate_MissingType(t *testing.T) {
	yaml := `
checks:
  - name: Go
    command: go version
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestValidate_WithTimeout(t *testing.T) {
	yaml := `
checks:
  - name: Slow
    type: command
    command: sleep 1
    timeout: 5s
`
	cfg, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cfg.Checks[0].Timeout != "5s" {
		t.Errorf("expected timeout '5s', got %q", cfg.Checks[0].Timeout)
	}
}

func TestValidate_WithGroups(t *testing.T) {
	yaml := `
checks:
  - name: DB
    type: port
    port: 5432
    groups: [backend, test]
`
	cfg, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(cfg.Checks[0].Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(cfg.Checks[0].Groups))
	}
}

func TestValidate_WithExtends(t *testing.T) {
	yaml := `
extends: base.pulse.yaml
checks:
  - name: Go
    type: command
    command: go version
`
	cfg, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if cfg.Extends != "base.pulse.yaml" {
		t.Errorf("expected extends 'base.pulse.yaml', got %q", cfg.Extends)
	}
}
