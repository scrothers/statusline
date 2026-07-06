package theme

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/scrothers/statusline/internal/style"
)

//go:embed themes/*.toml
var themeFS embed.FS

// DefaultName is the theme used when none is configured, or when a
// configured name doesn't match a built-in theme. claude-dark is a
// reasonable default for the widest audience with zero config: most
// terminals are dark, and it doesn't assume familiarity with any specific
// developer color scheme the way gruvbox/nord/dracula/etc. do.
const DefaultName = "claude-dark"

// names is the built-in theme list in a fixed display order (matching the
// README's theme table), since a map (as LoadRegistry returns) has no
// stable iteration order. DefaultName is listed first.
var names = []string{"claude-dark", "claude-light", "gruvbox", "catppuccin-mocha", "tokyo-night", "nord", "dracula"}

// Names returns the built-in theme names in a fixed, stable display order.
func Names() []string {
	out := make([]string, len(names))
	copy(out, names)
	return out
}

// Theme is the resolved color-token palette segments render against. Every
// theme fills the same fields, so rendering code never branches on theme
// name — it reads roles like Success/Warning/Danger uniformly. There is
// deliberately no background token: the statusline never paints a
// background, only foreground text and icons, so a theme is a foreground
// palette only.
type Theme struct {
	Name                         string
	IdentityAccent, IdentityText style.Color
	TextPrimary, TextSecondary   style.Color
	Success, Warning, Danger     style.Color
	Info, Muted                  style.Color
	TrackDim                     style.Color
}

type rawTheme struct {
	Name   string `toml:"name"`
	Colors struct {
		IdentityAccent string `toml:"identity_accent"`
		IdentityText   string `toml:"identity_text"`
		TextPrimary    string `toml:"text_primary"`
		TextSecondary  string `toml:"text_secondary"`
		Success        string `toml:"success"`
		Warning        string `toml:"warning"`
		Danger         string `toml:"danger"`
		Info           string `toml:"info"`
		Muted          string `toml:"muted"`
		TrackDim       string `toml:"track_dim"`
	} `toml:"colors"`
}

// LoadRegistry parses every embedded built-in theme file into a Theme,
// keyed by theme name.
func LoadRegistry() (map[string]Theme, error) {
	entries, err := fs.ReadDir(themeFS, "themes")
	if err != nil {
		return nil, fmt.Errorf("theme: read embedded themes dir: %w", err)
	}

	registry := make(map[string]Theme, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}
		data, err := themeFS.ReadFile("themes/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("theme: read %s: %w", entry.Name(), err)
		}
		var raw rawTheme
		if _, err := toml.Decode(string(data), &raw); err != nil {
			return nil, fmt.Errorf("theme: decode %s: %w", entry.Name(), err)
		}
		th, err := raw.resolve()
		if err != nil {
			return nil, err
		}
		registry[th.Name] = th
	}
	return registry, nil
}

func (r rawTheme) resolve() (Theme, error) {
	var errs []error
	parse := func(hex string) style.Color {
		c, err := style.ParseHex(hex)
		if err != nil {
			errs = append(errs, err)
		}
		return c
	}

	t := Theme{
		Name:           r.Name,
		IdentityAccent: parse(r.Colors.IdentityAccent),
		IdentityText:   parse(r.Colors.IdentityText),
		TextPrimary:    parse(r.Colors.TextPrimary),
		TextSecondary:  parse(r.Colors.TextSecondary),
		Success:        parse(r.Colors.Success),
		Warning:        parse(r.Colors.Warning),
		Danger:         parse(r.Colors.Danger),
		Info:           parse(r.Colors.Info),
		Muted:          parse(r.Colors.Muted),
		TrackDim:       parse(r.Colors.TrackDim),
	}
	if len(errs) > 0 {
		return Theme{}, fmt.Errorf("theme: resolve %q: %w", r.Name, errors.Join(errs...))
	}
	return t, nil
}

// Resolve looks up name in registry, falling back to DefaultName (and
// reporting that fallback via the second return value) when name is empty
// or unrecognized.
func Resolve(registry map[string]Theme, name string) (Theme, string) {
	if name == "" {
		name = DefaultName
	}
	if th, ok := registry[name]; ok {
		return th, ""
	}
	warning := fmt.Sprintf("statusline: unknown theme %q, using %q", name, DefaultName)
	return registry[DefaultName], warning
}

// tokenSetters maps a config-facing token name to the field it overrides.
// Pointer receiver is required: the returned pointers must address t's own
// fields, not a value-receiver copy's.
func (t *Theme) tokenSetters() map[string]*style.Color {
	return map[string]*style.Color{
		"identity_accent": &t.IdentityAccent,
		"identity_text":   &t.IdentityText,
		"text_primary":    &t.TextPrimary,
		"text_secondary":  &t.TextSecondary,
		"success":         &t.Success,
		"warning":         &t.Warning,
		"danger":          &t.Danger,
		"info":            &t.Info,
		"muted":           &t.Muted,
		"track_dim":       &t.TrackDim,
	}
}

// WithOverrides returns a copy of t with any recognized token names in
// overrides (hex color strings, e.g. {"success": "#00ff00"}) applied.
// Unrecognized token names or unparseable hex values are collected and
// returned as a single joined error; the returned Theme still has every
// other override applied.
func (t Theme) WithOverrides(overrides map[string]string) (Theme, error) {
	if len(overrides) == 0 {
		return t, nil
	}
	out := t
	setters := out.tokenSetters()
	var errs []error
	for name, hex := range overrides {
		field, ok := setters[name]
		if !ok {
			errs = append(errs, fmt.Errorf("theme: unknown color token %q", name))
			continue
		}
		c, err := style.ParseHex(hex)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		*field = c
	}
	if len(errs) > 0 {
		return out, fmt.Errorf("theme: apply overrides: %w", errors.Join(errs...))
	}
	return out, nil
}
