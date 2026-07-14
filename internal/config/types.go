package config

// Config is the fully-resolved statusline configuration: the built-in
// default with any user overlay merged on top.
type Config struct {
	Theme          string                   `toml:"theme"`
	ThemeOverrides map[string]string        `toml:"theme_overrides"`
	NerdFont       *bool                    `toml:"nerd_font"`
	Lines          []LineConfig             `toml:"lines"`
	Segments       map[string]SegmentConfig `toml:"segments"`
	Git            GitConfig                `toml:"git"`
	Budget         BudgetConfig             `toml:"budget"`
	// Provider forces the provider badge segment to a specific gateway
	// ("aws", "gcp", "azure", or "router") instead of auto-detecting from
	// the model id. This is the only way the Azure badge can ever appear —
	// Microsoft Foundry (and especially Azure AI Foundry deployment names)
	// carries no reliable id-shape signal to detect automatically. An
	// unrecognized value is treated the same as unset (auto-detect).
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
