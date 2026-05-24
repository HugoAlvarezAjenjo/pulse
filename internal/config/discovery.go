package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// configFileNames defines the config file names to search for, in priority order.
var configFileNames = []string{
	".pulse.yaml",
	"pulse.yaml",
}

// localConfigName is the optional local override file.
const localConfigName = ".pulse.local.yaml"

// Discover searches for a pulse config file in the given directory.
// Returns the path to the first found config file.
func Discover(dir string) (string, error) {
	for _, name := range configFileNames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no pulse config found (looked for %v in %s)", configFileNames, dir)
}

// DiscoverLocal checks if a local override config exists in the given directory.
// Returns the path if found, empty string otherwise.
func DiscoverLocal(dir string) string {
	path := filepath.Join(dir, localConfigName)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	return ""
}

// DiscoverAndLoad finds and loads a pulse config from the given directory.
// Resolves `extends` inheritance, then merges .pulse.local.yaml on top.
func DiscoverAndLoad(dir string) (*Config, string, error) {
	path, err := Discover(dir)
	if err != nil {
		return nil, "", err
	}

	// Load with extends resolution
	cfg, err := LoadWithExtends(path)
	if err != nil {
		return nil, path, err
	}

	// Try to load and merge local overrides
	if localPath := DiscoverLocal(dir); localPath != "" {
		localCfg, err := loadLocalConfig(localPath)
		if err == nil && localCfg != nil {
			cfg = Merge(cfg, localCfg)
		}
	}

	return cfg, path, nil
}

// loadLocalConfig loads a local config file, allowing empty checks (no validation).
func loadLocalConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse without strict validation (local can have partial config)
	return parseLocal(data)
}
