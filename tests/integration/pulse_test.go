package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	name := "pulse"
	if runtime.GOOS == "windows" {
		name = "pulse.exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", binary, "github.com/HugoAlvarezAjenjo/pulse/cmd/pulse")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build pulse: %v\n%s", err, out)
	}
	return binary
}

func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".pulse.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestIntegration_AllPassingChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Echo
    type: command
    command: echo hello

  - name: Config exists
    type: file
    path: .pulse.yaml
`)

	cmd := exec.Command(binary, "--plain")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "PASS") {
		t.Errorf("expected PASS in output:\n%s", output)
	}
	if !strings.Contains(output, "2 passed") {
		t.Errorf("expected '2 passed' in summary:\n%s", output)
	}
}

func TestIntegration_FailingCheck_ExitCode1(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Missing file
    type: file
    path: /nonexistent/file.txt
`)

	cmd := exec.Command(binary, "--plain")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected non-zero exit code")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d\n%s", exitErr.ExitCode(), out)
	}
}

func TestIntegration_JSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Echo
    type: command
    command: echo v1.0.0
    expected: ">= 1.0.0"
`)

	cmd := exec.Command(binary, "-o", "json")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0: %v\n%s", err, out)
	}

	var result struct {
		Results []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"results"`
		Summary struct {
			Total  int `json:"total"`
			Passed int `json:"passed"`
		} `json:"summary"`
	}

	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, out)
	}

	if result.Summary.Total != 1 {
		t.Errorf("expected total=1, got %d", result.Summary.Total)
	}
	if result.Summary.Passed != 1 {
		t.Errorf("expected passed=1, got %d", result.Summary.Passed)
	}
}

func TestIntegration_InvalidConfig_ExitCode2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks: []`)

	cmd := exec.Command(binary)
	cmd.Dir = dir
	_, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected non-zero exit code for invalid config")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
	}
}

func TestIntegration_Groups(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Frontend
    type: command
    command: echo frontend
    groups: [frontend]

  - name: Backend
    type: command
    command: echo backend
    groups: [backend]

  - name: Global
    type: command
    command: echo global
`)

	cmd := exec.Command(binary, "--plain", "-g", "frontend")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Frontend") {
		t.Error("expected Frontend check in output")
	}
	if !strings.Contains(output, "Global") {
		t.Error("expected Global check (no group = always runs)")
	}
	if strings.Contains(output, "Backend") {
		t.Error("Backend check should be filtered out")
	}
}

func TestIntegration_Validate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Test
    type: command
    command: echo ok
`)

	cmd := exec.Command(binary, "validate")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("validate should succeed: %v\n%s", err, out)
	}
}

func TestIntegration_List(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	binary := buildBinary(t)
	dir := t.TempDir()

	writeConfig(t, dir, `checks:
  - name: Alpha
    type: command
    command: echo a

  - name: Beta
    type: file
    path: /tmp
`)

	cmd := exec.Command(binary, "list")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("list should succeed: %v\n%s", err, out)
	}

	output := string(out)
	if !strings.Contains(output, "Alpha") || !strings.Contains(output, "Beta") {
		t.Errorf("list should show both checks:\n%s", output)
	}
}
