package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Load resolves configuration by merging the built-in default with an
// optional user file. explicitPath (from a --config flag) takes precedence
// over the fixed discovery order; pass "" to use discovery only. Load never
// fails: any read or parse problem is reported as a warning string and the
// built-in default (or as much of the merge as succeeded) is returned
// instead, since stdout is reserved for the rendered status line and a
// broken config must never blank it out.
func Load(explicitPath string) (Config, []string) {
	base := Default()
	var warnings []string

	path := explicitPath
	if path != "" {
		if _, err := os.Stat(path); err != nil {
			warnings = append(warnings, fmt.Sprintf("statusline: config %s: %v (using defaults)", path, err))
			return base, warnings
		}
	} else {
		path = firstExisting(candidatePaths())
		if path == "" {
			return base, warnings
		}
	}

	// #nosec G304 -- path is either the user's own --config flag or one of
	// the two fixed, well-known config locations in candidatePaths(); both
	// are the same trust level as any CLI tool's config file argument.
	data, err := os.ReadFile(path)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("statusline: reading config %s: %v (using defaults)", path, err))
		return base, warnings
	}

	var overlay Config
	if _, err := toml.Decode(string(data), &overlay); err != nil {
		warnings = append(warnings, fmt.Sprintf("statusline: parsing config %s: %v (using defaults)", path, err))
		return base, warnings
	}

	mergeInto(&base, overlay)
	return base, warnings
}

// candidatePaths lists non-explicit config locations in precedence order:
// $XDG_CONFIG_HOME/statusline/config.toml (or ~/.config/... if XDG_CONFIG_HOME
// is unset), then ~/.claude/statusline-config.toml.
func candidatePaths() []string {
	var paths []string

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "statusline", "config.toml"))
	} else if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "statusline", "config.toml"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".claude", "statusline-config.toml"))
	}
	return paths
}

func firstExisting(paths []string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
