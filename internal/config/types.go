package config

// Config is the fully-resolved statusline configuration: the built-in
// default with any user overlay merged on top.
type Config struct {
	// Theme manually forces the base palette to "dark" or "light" instead
	// of auto-detecting it from Claude Code's own settings.json (see
	// internal/claudetheme). Empty (the common case) or any other value
	// means "auto-detect"; this is a narrow escape hatch for previewing or
	// running the binary outside Claude Code, not a return to picking
	// between named aesthetic themes.
	Theme          string                   `toml:"theme"`
	ThemeOverrides map[string]string        `toml:"theme_overrides"`
	NerdFont       *bool                    `toml:"nerd_font"`
	Lines          []LineConfig             `toml:"lines"`
	Segments       map[string]SegmentConfig `toml:"segments"`
	Git            GitConfig                `toml:"git"`
	Budget         BudgetConfig             `toml:"budget"`
	// Provider forces the provider badge segment to a specific gateway
	// ("aws", "gcp", "azure", "cloudflare", "digitalocean", "router", or
	// "gateway") instead of auto-detecting. Auto-detection itself has two
	// tiers: Claude Code's own routing environment variables (see
	// DetectProviderFromEnv) are checked first — the only reliable way to
	// detect Azure/Foundry, Cloudflare, DigitalOcean, or a bare
	// corporate-relayed id, since none of them carry any distinguishing
	// shape in the model id — then the model id's own shape as a last
	// resort. This field is a manual override on top of both, not the only
	// way any particular badge can appear. An unrecognized value is treated
	// the same as unset (fall through to auto-detection).
	Provider string `toml:"provider"`
}

// NerdFontEnabled reports whether Nerd Font glyphs should render, defaulting
// to true when unset.
func (c Config) NerdFontEnabled() bool {
	return c.NerdFont == nil || *c.NerdFont
}

// LineConfig lists the segment IDs, in order, rendered on one output line.
type LineConfig struct {
	Enabled  bool     `toml:"enabled"`
	Segments []string `toml:"segments"`
}

// SegmentConfig carries a per-segment override. A nil Enabled means
// "unset" (segment decides for itself); an empty Icon means "use the
// built-in glyph".
type SegmentConfig struct {
	Enabled *bool  `toml:"enabled"`
	Icon    string `toml:"icon"`
}

// GitConfig tunes git status collection.
type GitConfig struct {
	Enabled    *bool `toml:"enabled"`
	TimeoutMS  int   `toml:"timeout_ms"`
	CacheTTLMS int   `toml:"cache_ttl_ms"`
}

// IsEnabled reports whether the git segment should collect status,
// defaulting to true when unset.
func (g GitConfig) IsEnabled() bool {
	return g.Enabled == nil || *g.Enabled
}

// BudgetConfig bounds the total time the binary may spend rendering.
type BudgetConfig struct {
	TotalTimeoutMS int `toml:"total_timeout_ms"`
}
