package config

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateHome points HOME (and clears XDG_CONFIG_HOME) at a fresh temp dir
// so candidatePaths() resolves deterministically. Mutating process env
// means the calling test cannot run in parallel with others that do the
// same.
func isolateHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	return home
}

func TestLoad_noConfigAnywhereUsesDefaults(t *testing.T) {
	isolateHome(t)

	cfg, warnings := Load("")
	if len(warnings) != 0 {
		t.Errorf("Load() warnings = %v, want none", warnings)
	}
	if cfg.Theme != "gruvbox" {
		t.Errorf("Theme = %q, want gruvbox", cfg.Theme)
	}
	if len(cfg.Lines) != 3 {
		t.Errorf("len(Lines) = %d, want 3", len(cfg.Lines))
	}
}

func TestLoad_explicitPathMissingWarnsAndFallsBack(t *testing.T) {
	isolateHome(t)

	cfg, warnings := Load("/nonexistent/path/config.toml")
	if len(warnings) != 1 {
		t.Fatalf("Load() warnings = %v, want exactly 1", warnings)
	}
	if cfg.Theme != "gruvbox" {
		t.Errorf("Theme = %q, want gruvbox (default)", cfg.Theme)
	}
}

func TestLoad_explicitPathOverridesScalarAndMergesMaps(t *testing.T) {
	isolateHome(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	body := `
theme = "nord"

[theme_overrides]
success = "#00ff00"

[segments.pr]
enabled = false
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, warnings := Load(path)
	if len(warnings) != 0 {
		t.Fatalf("Load() warnings = %v, want none", warnings)
	}
	if cfg.Theme != "nord" {
		t.Errorf("Theme = %q, want nord", cfg.Theme)
	}
	if cfg.ThemeOverrides["success"] != "#00ff00" {
		t.Errorf("ThemeOverrides[success] = %q, want #00ff00", cfg.ThemeOverrides["success"])
	}
	prCfg, ok := cfg.Segments["pr"]
	if !ok || prCfg.Enabled == nil || *prCfg.Enabled {
		t.Errorf("Segments[pr] = %+v, want Enabled=false", prCfg)
	}
	// Untouched defaults must survive the merge.
	if !cfg.NerdFontEnabled() {
		t.Error("NerdFontEnabled() = false, want true (untouched default)")
	}
	if len(cfg.Lines) != 3 {
		t.Errorf("len(Lines) = %d, want 3 (untouched default)", len(cfg.Lines))
	}
}

func TestLoad_malformedConfigWarnsAndFallsBack(t *testing.T) {
	isolateHome(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("not = [valid toml"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, warnings := Load(path)
	if len(warnings) != 1 {
		t.Fatalf("Load() warnings = %v, want exactly 1", warnings)
	}
	if cfg.Theme != "gruvbox" {
		t.Errorf("Theme = %q, want gruvbox (default)", cfg.Theme)
	}
}

func TestLoad_xdgConfigHomeDiscovered(t *testing.T) {
	home := isolateHome(t)
	xdg := filepath.Join(home, "xdgconf")
	if err := os.MkdirAll(filepath.Join(xdg, "statusline"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	t.Setenv("XDG_CONFIG_HOME", xdg)

	path := filepath.Join(xdg, "statusline", "config.toml")
	if err := os.WriteFile(path, []byte(`theme = "dracula"`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, warnings := Load("")
	if len(warnings) != 0 {
		t.Errorf("Load() warnings = %v, want none", warnings)
	}
	if cfg.Theme != "dracula" {
		t.Errorf("Theme = %q, want dracula", cfg.Theme)
	}
}

func TestLoad_claudeConfigDiscoveredWhenNoXDG(t *testing.T) {
	home := isolateHome(t)
	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	path := filepath.Join(claudeDir, "statusline-config.toml")
	if err := os.WriteFile(path, []byte(`theme = "tokyo-night"`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, warnings := Load("")
	if len(warnings) != 0 {
		t.Errorf("Load() warnings = %v, want none", warnings)
	}
	if cfg.Theme != "tokyo-night" {
		t.Errorf("Theme = %q, want tokyo-night", cfg.Theme)
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()
	cfg := Default()
	if cfg.Theme != "gruvbox" {
		t.Errorf("Default().Theme = %q, want gruvbox", cfg.Theme)
	}
	if !cfg.NerdFontEnabled() {
		t.Error("Default().NerdFontEnabled() = false, want true")
	}
	if !cfg.Git.IsEnabled() {
		t.Error("Default().Git.IsEnabled() = false, want true")
	}
	if cfg.Git.TimeoutMS == 0 || cfg.Git.CacheTTLMS == 0 || cfg.Budget.TotalTimeoutMS == 0 {
		t.Errorf("Default() has an unset timing field: %+v", cfg)
	}
}
