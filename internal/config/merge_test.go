package config

import (
	"testing"
)

func TestMerge_Override(t *testing.T) {
	base := &Config{
		Checks: []CheckConfig{
			{Name: "PostgreSQL", Type: "port", Host: "localhost", Port: 5432},
			{Name: "Redis", Type: "port", Host: "localhost", Port: 6379},
		},
	}

	local := &Config{
		Checks: []CheckConfig{
			{Name: "PostgreSQL", Type: "port", Host: "localhost", Port: 5433},
		},
	}

	merged := Merge(base, local)

	if len(merged.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(merged.Checks))
	}

	// PostgreSQL should be overridden
	if merged.Checks[0].Port != 5433 {
		t.Errorf("expected port 5433, got %d", merged.Checks[0].Port)
	}

	// Redis should be kept
	if merged.Checks[1].Name != "Redis" {
		t.Errorf("expected Redis, got %s", merged.Checks[1].Name)
	}
}

func TestMerge_Append(t *testing.T) {
	base := &Config{
		Checks: []CheckConfig{
			{Name: "Go", Type: "command", Command: "go version"},
		},
	}

	local := &Config{
		Checks: []CheckConfig{
			{Name: "Debug env", Type: "env", Variable: "DEBUG"},
		},
	}

	merged := Merge(base, local)

	if len(merged.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(merged.Checks))
	}

	if merged.Checks[0].Name != "Go" {
		t.Errorf("expected Go first, got %s", merged.Checks[0].Name)
	}
	if merged.Checks[1].Name != "Debug env" {
		t.Errorf("expected Debug env second, got %s", merged.Checks[1].Name)
	}
}

func TestMerge_NilLocal(t *testing.T) {
	base := &Config{
		Checks: []CheckConfig{
			{Name: "Go", Type: "command"},
		},
	}

	merged := Merge(base, nil)

	if len(merged.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(merged.Checks))
	}
}

func TestMerge_EmptyLocal(t *testing.T) {
	base := &Config{
		Checks: []CheckConfig{
			{Name: "Go", Type: "command"},
		},
	}

	local := &Config{Checks: []CheckConfig{}}

	merged := Merge(base, local)

	if len(merged.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(merged.Checks))
	}
}

func TestMerge_OverrideAndAppend(t *testing.T) {
	base := &Config{
		Checks: []CheckConfig{
			{Name: "Node.js", Type: "command", Command: "node --version", Expected: ">=18"},
			{Name: "DB", Type: "port", Port: 5432},
		},
	}

	local := &Config{
		Checks: []CheckConfig{
			{Name: "DB", Type: "port", Port: 5433},           // override
			{Name: "Extra", Type: "env", Variable: "SECRET"}, // new
		},
	}

	merged := Merge(base, local)

	if len(merged.Checks) != 3 {
		t.Fatalf("expected 3 checks, got %d", len(merged.Checks))
	}

	if merged.Checks[0].Name != "Node.js" {
		t.Errorf("expected Node.js first")
	}
	if merged.Checks[1].Name != "DB" || merged.Checks[1].Port != 5433 {
		t.Errorf("expected DB with port 5433, got port %d", merged.Checks[1].Port)
	}
	if merged.Checks[2].Name != "Extra" {
		t.Errorf("expected Extra third")
	}
}
