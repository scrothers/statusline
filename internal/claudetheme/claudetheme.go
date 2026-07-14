package claudetheme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Resolve returns "dark" or "light": Claude Code's own effective theme,
// collapsed to a base. It reads Claude Code's settings.json files in the
// same scope order Claude Code itself uses — local project overrides
// project overrides user — falling through to Claude Code's own documented
// default ("dark") when no scope sets a theme, and to "dark" again on any
// unreadable/unrecognized value, since a detection problem must never block
// rendering. projectDir is typically payload.Workspace.ProjectDir; an empty
// value just skips the two project-scoped lookups.
func Resolve(projectDir string) (string, []string) {
	raw := "dark"
	for _, path := range settingsPaths(projectDir) {
		if theme, ok := settingsTheme(path); ok {
			raw = theme
			break
		}
	}
	return classify(raw)
}

// settingsPaths lists Claude Code's own settings.json locations in
// precedence order (first hit wins): local project settings, project
// settings, then user settings — mirroring the "Local > Project > User"
// scope order from Claude Code's own settings precedence (the "Managed"
// and "CLI argument" scopes above these aren't file-based and aren't
// visible to a separate statusline process).
func settingsPaths(projectDir string) []string {
	var paths []string
	if projectDir != "" {
		paths = append(paths,
			filepath.Join(projectDir, ".claude", "settings.local.json"),
			filepath.Join(projectDir, ".claude", "settings.json"),
		)
	}
	paths = append(paths, filepath.Join(configDir(), "settings.json"))
	return paths
}

// configDir returns Claude Code's own config directory: $CLAUDE_CONFIG_DIR
// if set, otherwise ~/.claude.
func configDir() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude")
}

// settingsTheme reads one settings.json file and returns its "theme" key.
// Every other key is ignored — statusline has no need to model Claude
// Code's full settings schema. A missing file, a read error, malformed
// JSON, or an empty/absent theme key all report ok=false, so the caller
// falls through to the next scope.
func settingsTheme(path string) (string, bool) {
	// #nosec G304 -- path is one of a fixed set of well-known Claude Code
	// settings locations, the same trust level as internal/config's own
	// config file discovery.
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	var parsed struct {
		Theme string `json:"theme"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", false
	}
	if parsed.Theme == "" {
		return "", false
	}
	return parsed.Theme, true
}

// classify maps a raw Claude Code theme value to "dark" or "light",
// resolving "auto" and "custom:<slug>" references as needed. It always
// returns a base and never fails — an unrecognized value falls back to
// "dark" with an explanatory warning rather than blocking rendering.
func classify(raw string) (string, []string) {
	switch raw {
	case "dark", "dark-daltonized", "dark-ansi":
		return "dark", nil
	case "light", "light-daltonized", "light-ansi":
		return "light", nil
	case "auto":
		return resolveAuto()
	}

	if slug, ok := strings.CutPrefix(raw, "custom:"); ok {
		return classifyCustom(slug)
	}

	return "dark", []string{
		fmt.Sprintf("statusline: unrecognized Claude Code theme %q, using \"dark\"", raw),
	}
}

// classifyCustom resolves a "custom:<slug>" reference by reading its base
// preset from ~/.claude/themes/<slug>.json. Custom themes can also override
// individual color tokens, but statusline doesn't parse those (see the
// package doc / project plan) — only the base preset is honored. A
// "custom:<plugin>:<slug>" reference (plugin-contributed, not a file under
// ~/.claude/themes/) isn't resolvable from a plain settings read, so it
// falls back to "dark" like any other unreadable custom theme.
func classifyCustom(slug string) (string, []string) {
	if slug == "" || strings.Contains(slug, ":") {
		return "dark", []string{
			fmt.Sprintf("statusline: unsupported custom Claude Code theme %q, using \"dark\"", "custom:"+slug),
		}
	}

	path := filepath.Join(configDir(), "themes", slug+".json")
	// #nosec G304 -- path is built from a fixed directory plus a slug
	// already validated to contain no path separators or colons.
	data, err := os.ReadFile(path)
	if err != nil {
		return "dark", []string{
			fmt.Sprintf("statusline: reading custom Claude Code theme %q: %v, using \"dark\"", slug, err),
		}
	}
	var parsed struct {
		Base string `json:"base"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "dark", []string{
			fmt.Sprintf("statusline: parsing custom Claude Code theme %q: %v, using \"dark\"", slug, err),
		}
	}
	if parsed.Base == "" {
		return "dark", nil // Claude Code's own documented default base
	}
	// parsed.Base is one of the six built-in presets by Claude Code's own
	// spec, never "auto" or another "custom:" reference, so this can't
	// recurse more than once.
	return classify(parsed.Base)
}

// resolveAuto best-effort approximates Claude Code's "auto" (terminal
// light/dark background detection) via the COLORFGBG environment
// variable, which many terminals and tmux configs set — never by
// querying the terminal directly, since statusline's stdout is piped to
// Claude Code rather than connected to the terminal. Absent, malformed,
// or ambiguous values all fall back to "dark".
func resolveAuto() (string, []string) {
	raw := os.Getenv("COLORFGBG")
	if raw == "" {
		return "dark", nil
	}

	// Format is typically "fg;bg", though some terminals emit a third
	// field ("fg;bogus;bg"); the background is always the last field.
	fields := strings.Split(raw, ";")
	bg, err := strconv.Atoi(strings.TrimSpace(fields[len(fields)-1]))
	if err != nil {
		return "dark", nil
	}

	// ANSI background 7 (white) or 15 (bright white) reads as a light
	// terminal background; every other index, including the 8 darker/
	// non-white colors, reads as dark.
	if bg == 7 || bg == 15 {
		return "light", nil
	}
	return "dark", nil
}
