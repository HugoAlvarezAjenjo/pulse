package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover_PulseYaml(t *testing.T) {
	dir := t.TempDir()

	// Create .pulse.yaml
	path := filepath.Join(dir, ".pulse.yaml")
	if err := os.WriteFile(path, []byte("checks:\n  - name: test\n    type: command\n    command: echo hi\n"), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := Discover(dir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if found != path {
		t.Errorf("expected %s, got %s", path, found)
	}
}

func TestDiscover_FallbackPulseYaml(t *testing.T) {
	dir := t.TempDir()

	// Create pulse.yaml (without dot prefix)
	path := filepath.Join(dir, "pulse.yaml")
	if err := os.WriteFile(path, []byte("checks:\n  - name: test\n    type: command\n    command: echo hi\n"), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := Discover(dir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if found != path {
		t.Errorf("expected %s, got %s", path, found)
	}
}

func TestDiscover_Priority(t *testing.T) {
	dir := t.TempDir()

	// Create both files - .pulse.yaml should win
	dotPath := filepath.Join(dir, ".pulse.yaml")
	plainPath := filepath.Join(dir, "pulse.yaml")

	if err := os.WriteFile(dotPath, []byte("checks:\n  - name: test\n    type: command\n    command: echo hi\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(plainPath, []byte("checks:\n  - name: test\n    type: command\n    command: echo hi\n"), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := Discover(dir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if found != dotPath {
		t.Errorf("expected .pulse.yaml to take priority, got %s", found)
	}
}

func TestDiscover_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := Discover(dir)
	if err == nil {
		t.Fatal("expected error when no config found")
	}
}
