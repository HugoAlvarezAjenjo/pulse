package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level pulse configuration.
type Config struct {
	Checks []CheckConfig `yaml:"checks"`
}

// CheckConfig represents a single check entry in the configuration.
type CheckConfig struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Command   string     `yaml:"command,omitempty"`
	Expected  string     `yaml:"expected,omitempty"`
	Path      string     `yaml:"path,omitempty"`
	Host      string     `yaml:"host,omitempty"`
	Port      int        `yaml:"port,omitempty"`
	Variable  string     `yaml:"variable,omitempty"`
	URL       string     `yaml:"url,omitempty"`
	Status    int        `yaml:"status,omitempty"`
	Container string     `yaml:"container,omitempty"`
	Fix       *FixConfig `yaml:"fix,omitempty"`
}

// FixConfig represents a fix command associated with a check.
type FixConfig struct {
	Run string `yaml:"run"`
}

// Load reads and parses a config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	return Parse(data)
}

// Parse parses raw YAML bytes into a Config.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate performs basic validation on the parsed config.
func validate(cfg *Config) error {
	if len(cfg.Checks) == 0 {
		return fmt.Errorf("config: no checks defined")
	}

	for i, check := range cfg.Checks {
		if check.Name == "" {
			return fmt.Errorf("config: check at index %d is missing 'name'", i)
		}
		if check.Type == "" {
			return fmt.Errorf("config: check %q is missing 'type'", check.Name)
		}
	}

	return nil
}
