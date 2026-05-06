package config

// Merge combines a base config with a local override config.
// Rules:
//   - Checks in local with the same name as base → override (replace completely)
//   - Checks in local with a new name → append
//   - Checks in base not in local → kept as-is
//   - Local cannot remove checks from base
func Merge(base, local *Config) *Config {
	if local == nil || len(local.Checks) == 0 {
		return base
	}

	// Index local checks by name for fast lookup
	localByName := make(map[string]CheckConfig, len(local.Checks))
	for _, c := range local.Checks {
		localByName[c.Name] = c
	}

	merged := &Config{
		Checks: make([]CheckConfig, 0, len(base.Checks)+len(local.Checks)),
	}

	// Process base checks: override if local has same name, keep otherwise
	seen := make(map[string]bool)
	for _, check := range base.Checks {
		if override, exists := localByName[check.Name]; exists {
			merged.Checks = append(merged.Checks, override)
			seen[check.Name] = true
		} else {
			merged.Checks = append(merged.Checks, check)
		}
	}

	// Append local-only checks (new checks not in base)
	for _, check := range local.Checks {
		if !seen[check.Name] {
			merged.Checks = append(merged.Checks, check)
		}
	}

	return merged
}
