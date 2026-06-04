package checks

import (
	"testing"

	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
)

func TestFromConfig_ValidTypes(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.CheckConfig
	}{
		{"command", config.CheckConfig{Name: "go", Type: "command", Command: "go version"}},
		{"file", config.CheckConfig{Name: "mod", Type: "file", Path: "go.mod"}},
		{"port", config.CheckConfig{Name: "db", Type: "port", Port: 5432}},
		{"env", config.CheckConfig{Name: "key", Type: "env", Variable: "HOME"}},
		{"http", config.CheckConfig{Name: "api", Type: "http", URL: "http://localhost/health"}},
		{"docker", config.CheckConfig{Name: "redis", Type: "docker", Container: "redis"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check, err := FromConfig(tt.cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if check == nil {
				t.Fatal("expected non-nil check")
			}
		})
	}
}

func TestFromConfig_UnknownType(t *testing.T) {
	cfg := config.CheckConfig{Name: "test", Type: "unknown"}
	_, err := FromConfig(cfg)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestFromConfig_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		cfg  config.CheckConfig
	}{
		{"command without command", config.CheckConfig{Name: "test", Type: "command"}},
		{"file without path", config.CheckConfig{Name: "test", Type: "file"}},
		{"port without port", config.CheckConfig{Name: "test", Type: "port"}},
		{"env without variable", config.CheckConfig{Name: "test", Type: "env"}},
		{"http without url", config.CheckConfig{Name: "test", Type: "http"}},
		{"docker without container", config.CheckConfig{Name: "test", Type: "docker"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromConfig(tt.cfg)
			if err == nil {
				t.Error("expected error for missing required field")
			}
		})
	}
}

func TestFromConfig_WithFix(t *testing.T) {
	cfg := config.CheckConfig{
		Name: "db",
		Type: "port",
		Port: 5432,
		Fix:  &config.FixConfig{Run: "docker compose up db -d"},
	}

	check, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pc, ok := check.(*PortCheck)
	if !ok {
		t.Fatal("expected *PortCheck")
	}
	if pc.Fix == nil {
		t.Error("expected fix to be set")
	}
	if pc.Fix.Run != "docker compose up db -d" {
		t.Errorf("fix.Run = %q, want %q", pc.Fix.Run, "docker compose up db -d")
	}
}

func TestFromConfig_PortDefaultHost(t *testing.T) {
	cfg := config.CheckConfig{
		Name: "db",
		Type: "port",
		Port: 5432,
	}

	check, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pc := check.(*PortCheck)
	if pc.Host != "localhost" {
		t.Errorf("expected default host 'localhost', got %q", pc.Host)
	}
}
