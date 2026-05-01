package checks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

func TestFileCheck_Exists(t *testing.T) {
	// Create a temp file
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	check := &FileCheck{
		Name: "test file",
		Path: path,
	}

	r := check.Run(context.Background())
	if r.Status != result.Success {
		t.Errorf("expected success, got %s: %s", r.Status, r.Message)
	}
}

func TestFileCheck_NotFound(t *testing.T) {
	check := &FileCheck{
		Name: "missing file",
		Path: "/nonexistent/path/to/file.txt",
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure, got %s", r.Status)
	}
}

func TestFileCheck_IsDirectory(t *testing.T) {
	dir := t.TempDir()

	check := &FileCheck{
		Name: "directory",
		Path: dir,
	}

	r := check.Run(context.Background())
	if r.Status != result.Failure {
		t.Errorf("expected failure for directory, got %s", r.Status)
	}
}
