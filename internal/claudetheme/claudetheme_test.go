package claudetheme

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateConfigDir points CLAUDE_CONFIG_DIR at a fresh temp dir so
// settingsPaths' user-scope lookup resolves deterministically, independent
// of any real ~/.claude on the machine running the test.
func isolateConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", dir)
	t.Setenv("COLORFGBG", "")
	return dir
}

func writeJSON(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func TestResolve_noSettingsAnywhereDefaultsToDark(t *testing.T) {
	isolateConfigDir(t)

	base, warnings := Resolve("")
	if base != "dark" {
		t.Errorf("base = %q, want dark", base)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}
}

func TestResolve_userScope(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "light"}`)

	base, _ := Resolve("")
	if base != "light" {
		t.Errorf("base = %q, want light", base)
	}
}

func TestResolve_projectOverridesUser(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "dark"}`)

	projectDir := t.TempDir()
	writeJSON(t, filepath.Join(projectDir, ".claude", "settings.json"), `{"theme": "light"}`)

	base, _ := Resolve(projectDir)
	if base != "light" {
		t.Errorf("base = %q, want light (project should override user)", base)
	}
}

func TestResolve_localOverridesProject(t *testing.T) {
	isolateConfigDir(t)

	projectDir := t.TempDir()
	writeJSON(t, filepath.Join(projectDir, ".claude", "settings.json"), `{"theme": "dark"}`)
	writeJSON(t, filepath.Join(projectDir, ".claude", "settings.local.json"), `{"theme": "light"}`)

	base, _ := Resolve(projectDir)
	if base != "light" {
		t.Errorf("base = %q, want light (local should override project)", base)
	}
}

func TestResolve_malformedTierFallsThroughToNext(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "light"}`)

	projectDir := t.TempDir()
	writeJSON(t, filepath.Join(projectDir, ".claude", "settings.json"), `not json at all`)

	base, _ := Resolve(projectDir)
	if base != "light" {
		t.Errorf("base = %q, want light (malformed project settings should fall through to user)", base)
	}
}

func TestResolve_emptyThemeKeyFallsThroughToNext(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "light"}`)

	projectDir := t.TempDir()
	writeJSON(t, filepath.Join(projectDir, ".claude", "settings.json"), `{"theme": ""}`)

	base, _ := Resolve(projectDir)
	if base != "light" {
		t.Errorf("base = %q, want light (empty project theme should fall through to user)", base)
	}
}

func TestClassify_presets(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{"dark", "dark"},
		{"dark-daltonized", "dark"},
		{"dark-ansi", "dark"},
		{"light", "light"},
		{"light-daltonized", "light"},
		{"light-ansi", "light"},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			isolateConfigDir(t)
			got, warnings := classify(tt.raw)
			if got != tt.want {
				t.Errorf("classify(%q) = %q, want %q", tt.raw, got, tt.want)
			}
			if len(warnings) != 0 {
				t.Errorf("classify(%q) warnings = %v, want none", tt.raw, warnings)
			}
		})
	}
}

func TestClassify_unrecognizedFallsBackToDarkWithWarning(t *testing.T) {
	base, warnings := classify("not-a-real-theme")
	if base != "dark" {
		t.Errorf("base = %q, want dark", base)
	}
	if len(warnings) != 1 {
		t.Errorf("warnings = %v, want exactly 1", warnings)
	}
}

func TestResolveAuto(t *testing.T) {
	tests := []struct {
		name      string
		colorfgbg string
		want      string
	}{
		{"unset", "", "dark"},
		{"black background", "15;0", "dark"},
		{"white background", "0;7", "light"},
		{"bright white background", "0;15", "light"},
		{"three field takes last as bg", "0;1;15", "light"},
		{"malformed", "not-a-number", "dark"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("COLORFGBG", tt.colorfgbg)
			got, _ := resolveAuto()
			if got != tt.want {
				t.Errorf("resolveAuto() with COLORFGBG=%q = %q, want %q", tt.colorfgbg, got, tt.want)
			}
		})
	}
}

func TestResolve_autoUsesColorfgbg(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "auto"}`)
	t.Setenv("COLORFGBG", "0;15")

	base, _ := Resolve("")
	if base != "light" {
		t.Errorf("base = %q, want light", base)
	}
}

func TestClassifyCustom_resolvesBaseFromThemeFile(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "themes", "midnight.json"), `{"name": "Midnight", "base": "light"}`)

	base, warnings := classifyCustom("midnight")
	if base != "light" {
		t.Errorf("base = %q, want light", base)
	}
	if len(warnings) != 0 {
		t.Errorf("warnings = %v, want none", warnings)
	}
}

func TestClassifyCustom_missingBaseDefaultsToDark(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "themes", "noop.json"), `{"name": "Noop"}`)

	base, _ := classifyCustom("noop")
	if base != "dark" {
		t.Errorf("base = %q, want dark", base)
	}
}

func TestClassifyCustom_unreadableFileFallsBackToDark(t *testing.T) {
	isolateConfigDir(t)

	base, warnings := classifyCustom("does-not-exist")
	if base != "dark" {
		t.Errorf("base = %q, want dark", base)
	}
	if len(warnings) != 1 {
		t.Errorf("warnings = %v, want exactly 1", warnings)
	}
}

func TestClassifyCustom_pluginQualifiedFallsBackToDark(t *testing.T) {
	isolateConfigDir(t)

	base, warnings := classifyCustom("some-plugin:midnight")
	if base != "dark" {
		t.Errorf("base = %q, want dark", base)
	}
	if len(warnings) != 1 {
		t.Errorf("warnings = %v, want exactly 1", warnings)
	}
}

func TestResolve_customThemeReference(t *testing.T) {
	configDir := isolateConfigDir(t)
	writeJSON(t, filepath.Join(configDir, "settings.json"), `{"theme": "custom:midnight"}`)
	writeJSON(t, filepath.Join(configDir, "themes", "midnight.json"), `{"base": "light"}`)

	base, _ := Resolve("")
	if base != "light" {
		t.Errorf("base = %q, want light", base)
	}
}
