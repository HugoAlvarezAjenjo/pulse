package generator

import (
	"strings"
	"testing"
)

func TestGetPreset_Known(t *testing.T) {
	knownPresets := []string{"go", "node", "python", "java", "rust"}

	for _, name := range knownPresets {
		tmpl, ok := GetPreset(name)
		if !ok {
			t.Errorf("GetPreset(%q) returned ok=false, want true", name)
			continue
		}
		if !strings.Contains(tmpl, "checks:") {
			t.Errorf("GetPreset(%q) template missing 'checks:' key", name)
		}
	}
}

func TestGetPreset_Unknown(t *testing.T) {
	_, ok := GetPreset("nonexistent-stack")
	if ok {
		t.Error("GetPreset(\"nonexistent-stack\") returned ok=true, want false")
	}
}

func TestAvailablePresets_NotEmpty(t *testing.T) {
	presets := AvailablePresets()
	if len(presets) == 0 {
		t.Error("AvailablePresets() returned empty list")
	}
}

func TestAvailablePresets_AllResolvable(t *testing.T) {
	for _, name := range AvailablePresets() {
		_, ok := GetPreset(name)
		if !ok {
			t.Errorf("AvailablePresets() listed %q but GetPreset returns ok=false", name)
		}
	}
}

func TestEmptyTemplate_Valid(t *testing.T) {
	if !strings.Contains(EmptyTemplate, "checks:") {
		t.Error("EmptyTemplate missing 'checks:' key")
	}
	if !strings.Contains(EmptyTemplate, "Example") {
		t.Error("EmptyTemplate missing example check")
	}
}
