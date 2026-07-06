package config

// mergeInto applies overlay on top of base, field by field. Zero/nil values
// in overlay mean "not set by the user" and leave base untouched; Lines
// replaces base wholesale when present at all (the user is taking over
// layout), while ThemeOverrides and Segments merge key by key so a user can
// override one entry without redefining every other one.
func mergeInto(base *Config, overlay Config) {
	if overlay.Theme != "" {
		base.Theme = overlay.Theme
	}
	for k, v := range overlay.ThemeOverrides {
		if base.ThemeOverrides == nil {
			base.ThemeOverrides = make(map[string]string, len(overlay.ThemeOverrides))
		}
		base.ThemeOverrides[k] = v
	}
	if overlay.NerdFont != nil {
		base.NerdFont = overlay.NerdFont
	}
	if len(overlay.Lines) > 0 {
		base.Lines = overlay.Lines
	}
	for id, seg := range overlay.Segments {
		existing := base.Segments[id]
		if seg.Enabled != nil {
			existing.Enabled = seg.Enabled
		}
		if seg.Icon != "" {
			existing.Icon = seg.Icon
		}
		if base.Segments == nil {
			base.Segments = make(map[string]SegmentConfig, len(overlay.Segments))
		}
		base.Segments[id] = existing
	}
	if overlay.Git.Enabled != nil {
		base.Git.Enabled = overlay.Git.Enabled
	}
	if overlay.Git.TimeoutMS != 0 {
		base.Git.TimeoutMS = overlay.Git.TimeoutMS
	}
	if overlay.Git.CacheTTLMS != 0 {
		base.Git.CacheTTLMS = overlay.Git.CacheTTLMS
	}
	if overlay.Budget.TotalTimeoutMS != 0 {
		base.Budget.TotalTimeoutMS = overlay.Budget.TotalTimeoutMS
	}
}
