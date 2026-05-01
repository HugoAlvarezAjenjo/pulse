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

// DiscoverAndLoad finds and loads a pulse config from the given directory.
func DiscoverAndLoad(dir string) (*Config, string, error) {
	path, err := Discover(dir)
	if err != nil {
		return nil, "", err
	}

	cfg, err := Load(path)
	if err != nil {
		return nil, path, err
	}

	return cfg, path, nil
}
