package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const maxExtendsDepth = 5

// LoadWithExtends loads a config file and resolves its `extends` chain.
// Supports relative and absolute paths. Prevents circular/deep inheritance.
func LoadWithExtends(path string) (*Config, error) {
	return loadWithDepth(path, 0, nil)
}

func loadWithDepth(path string, depth int, seen map[string]bool) (*Config, error) {
	if depth > maxExtendsDepth {
		return nil, fmt.Errorf("extends: max depth (%d) exceeded — possible circular inheritance", maxExtendsDepth)
	}

	if seen == nil {
		seen = make(map[string]bool)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("extends: resolving path %q: %w", path, err)
	}

	if seen[absPath] {
		return nil, fmt.Errorf("extends: circular inheritance detected at %q", path)
	}
	seen[absPath] = true

	// Load this config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}

	cfg, err := parseLocal(data) // Use parseLocal to allow extends-only files with no checks
	if err != nil {
		return nil, err
	}

	// If no extends, validate and return
	if cfg.Extends == "" {
		// Validate only leaf configs (or configs without extends)
		if len(cfg.Checks) == 0 {
			return nil, fmt.Errorf("config %q: no checks defined", path)
		}
		return cfg, nil
	}

	// Resolve extends path relative to the current config file
	parentPath := cfg.Extends
	if !filepath.IsAbs(parentPath) {
		parentPath = filepath.Join(filepath.Dir(path), parentPath)
	}

	// Load parent recursively
	parent, err := loadWithDepth(parentPath, depth+1, seen)
	if err != nil {
		return nil, fmt.Errorf("extends: loading %q from %q: %w", cfg.Extends, path, err)
	}

	// Merge: child overrides parent (same semantics as .pulse.local.yaml)
	merged := Merge(parent, cfg)
	return merged, nil
}
