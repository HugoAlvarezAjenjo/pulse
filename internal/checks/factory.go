package checks

import (
	"fmt"

	"github.com/HugoAlvarezAjenjo/pulse/internal/config"
	"github.com/HugoAlvarezAjenjo/pulse/internal/result"
)

// FromConfig creates a Check from a config.CheckConfig entry.
func FromConfig(cfg config.CheckConfig) (Check, error) {
	var fix *result.Fix
	if cfg.Fix != nil && cfg.Fix.Run != "" {
		fix = &result.Fix{Run: cfg.Fix.Run}
	}

	switch cfg.Type {
	case "command":
		if cfg.Command == "" {
			return nil, fmt.Errorf("check %q: command type requires 'command' field", cfg.Name)
		}
		return &CommandCheck{
			Name:     cfg.Name,
			Command:  cfg.Command,
			Expected: cfg.Expected,
			Fix:      fix,
		}, nil

	case "file":
		if cfg.Path == "" {
			return nil, fmt.Errorf("check %q: file type requires 'path' field", cfg.Name)
		}
		return &FileCheck{
			Name: cfg.Name,
			Path: cfg.Path,
			Fix:  fix,
		}, nil

	case "port":
		if cfg.Port == 0 {
			return nil, fmt.Errorf("check %q: port type requires 'port' field", cfg.Name)
		}
		host := cfg.Host
		if host == "" {
			host = "localhost"
		}
		return &PortCheck{
			Name: cfg.Name,
			Host: host,
			Port: cfg.Port,
			Fix:  fix,
		}, nil

	default:
		return nil, fmt.Errorf("check %q: unknown type %q", cfg.Name, cfg.Type)
	}
}
