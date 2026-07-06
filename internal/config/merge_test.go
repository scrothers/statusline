package config

import "testing"

func TestMergeInto(t *testing.T) {
	t.Parallel()

	t.Run("segment merge preserves untouched fields", func(t *testing.T) {
		t.Parallel()
		enabledFalse := false
		base := Config{
			Segments: map[string]SegmentConfig{
				"pr": {Icon: "custom-icon"},
			},
		}
		overlay := Config{
			Segments: map[string]SegmentConfig{
				"pr": {Enabled: &enabledFalse},
			},
		}
		mergeInto(&base, overlay)

		got := base.Segments["pr"]
		if got.Icon != "custom-icon" {
			t.Errorf("Icon = %q, want custom-icon (preserved)", got.Icon)
		}
		if got.Enabled == nil || *got.Enabled {
			t.Errorf("Enabled = %v, want false (overridden)", got.Enabled)
		}
	})

	t.Run("lines replace wholesale, not merge", func(t *testing.T) {
		t.Parallel()
		base := Config{Lines: []LineConfig{{Enabled: true, Segments: []string{"model"}}, {Enabled: true, Segments: []string{"git"}}}}
		overlay := Config{Lines: []LineConfig{{Enabled: true, Segments: []string{"only-this"}}}}
		mergeInto(&base, overlay)

		if len(base.Lines) != 1 || base.Lines[0].Segments[0] != "only-this" {
			t.Errorf("Lines = %+v, want wholesale replacement", base.Lines)
		}
	})

	t.Run("theme overrides merge key by key", func(t *testing.T) {
		t.Parallel()
		base := Config{ThemeOverrides: map[string]string{"success": "#111111"}}
		overlay := Config{ThemeOverrides: map[string]string{"danger": "#222222"}}
		mergeInto(&base, overlay)

		if base.ThemeOverrides["success"] != "#111111" {
			t.Errorf("success = %q, want preserved #111111", base.ThemeOverrides["success"])
		}
		if base.ThemeOverrides["danger"] != "#222222" {
			t.Errorf("danger = %q, want #222222", base.ThemeOverrides["danger"])
		}
	})

	t.Run("zero-value overlay changes nothing", func(t *testing.T) {
		t.Parallel()
		base := Default()
		want := Default()
		mergeInto(&base, Config{})
		if base.Theme != want.Theme || len(base.Lines) != len(want.Lines) {
			t.Errorf("mergeInto with zero-value overlay changed base: got %+v, want %+v", base, want)
		}
	})

	t.Run("segments map starts nil, an overlay entry initializes it", func(t *testing.T) {
		t.Parallel()
		base := Config{}
		overlay := Config{Segments: map[string]SegmentConfig{"pr": {Icon: "new-icon"}}}
		mergeInto(&base, overlay)

		if base.Segments == nil {
			t.Fatal("Segments = nil, want initialized")
		}
		if base.Segments["pr"].Icon != "new-icon" {
			t.Errorf("Icon = %q, want new-icon", base.Segments["pr"].Icon)
		}
	})

	t.Run("scalar overrides all apply", func(t *testing.T) {
		t.Parallel()
		nerdFontOff := false
		gitEnabledOff := false
		base := Default()
		overlay := Config{
			Theme:    "nord",
			NerdFont: &nerdFontOff,
			Git: GitConfig{
				Enabled:    &gitEnabledOff,
				TimeoutMS:  999,
				CacheTTLMS: 888,
			},
			Budget: BudgetConfig{TotalTimeoutMS: 777},
		}
		mergeInto(&base, overlay)

		if base.Theme != "nord" {
			t.Errorf("Theme = %q, want nord", base.Theme)
		}
		if base.NerdFontEnabled() {
			t.Error("NerdFontEnabled() = true, want false (overridden)")
		}
		if base.Git.IsEnabled() {
			t.Error("Git.IsEnabled() = true, want false (overridden)")
		}
		if base.Git.TimeoutMS != 999 {
			t.Errorf("Git.TimeoutMS = %d, want 999", base.Git.TimeoutMS)
		}
		if base.Git.CacheTTLMS != 888 {
			t.Errorf("Git.CacheTTLMS = %d, want 888", base.Git.CacheTTLMS)
		}
		if base.Budget.TotalTimeoutMS != 777 {
			t.Errorf("Budget.TotalTimeoutMS = %d, want 777", base.Budget.TotalTimeoutMS)
		}
	})
}
