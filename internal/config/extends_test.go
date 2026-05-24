package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithExtends_Simple(t *testing.T) {
	dir := t.TempDir()

	// Create parent
	parent := `checks:
  - name: Git
    type: command
    command: git --version
  - name: Docker
    type: command
    command: docker --version
`
	if err := os.WriteFile(filepath.Join(dir, "base.pulse.yaml"), []byte(parent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create child that extends parent
	child := `extends: base.pulse.yaml
checks:
  - name: Node.js
    type: command
    command: node --version
`
	childPath := filepath.Join(dir, ".pulse.yaml")
	if err := os.WriteFile(childPath, []byte(child), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadWithExtends(childPath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Should have 3 checks: Git + Docker from parent, Node.js from child
	if len(cfg.Checks) != 3 {
		t.Fatalf("expected 3 checks, got %d", len(cfg.Checks))
	}

	if cfg.Checks[0].Name != "Git" {
		t.Errorf("expected first check 'Git', got %q", cfg.Checks[0].Name)
	}
	if cfg.Checks[2].Name != "Node.js" {
		t.Errorf("expected third check 'Node.js', got %q", cfg.Checks[2].Name)
	}
}

func TestLoadWithExtends_Override(t *testing.T) {
	dir := t.TempDir()

	parent := `checks:
  - name: Database
    type: port
    host: localhost
    port: 5432
  - name: Git
    type: command
    command: git --version
`
	if err := os.WriteFile(filepath.Join(dir, "base.pulse.yaml"), []byte(parent), 0644); err != nil {
		t.Fatal(err)
	}

	child := `extends: base.pulse.yaml
checks:
  - name: Database
    type: port
    host: localhost
    port: 5433
`
	childPath := filepath.Join(dir, ".pulse.yaml")
	if err := os.WriteFile(childPath, []byte(child), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadWithExtends(childPath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Should have 2 checks: Database (overridden) + Git (inherited)
	if len(cfg.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(cfg.Checks))
	}

	if cfg.Checks[0].Port != 5433 {
		t.Errorf("expected overridden port 5433, got %d", cfg.Checks[0].Port)
	}
}

func TestLoadWithExtends_CircularDetection(t *testing.T) {
	dir := t.TempDir()

	// A extends B, B extends A
	a := `extends: b.yaml
checks:
  - name: A
    type: command
    command: echo a
`
	b := `extends: a.yaml
checks:
  - name: B
    type: command
    command: echo b
`
	if err := os.WriteFile(filepath.Join(dir, "a.yaml"), []byte(a), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.yaml"), []byte(b), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadWithExtends(filepath.Join(dir, "a.yaml"))
	if err == nil {
		t.Fatal("expected error for circular inheritance")
	}
}

func TestLoadWithExtends_NoExtends(t *testing.T) {
	dir := t.TempDir()

	content := `checks:
  - name: Go
    type: command
    command: go version
`
	path := filepath.Join(dir, ".pulse.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadWithExtends(path)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(cfg.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(cfg.Checks))
	}
}
