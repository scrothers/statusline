package theme

import "testing"

func TestGlyph(t *testing.T) {
	t.Parallel()

	t.Run("nerd font enabled returns glyph", func(t *testing.T) {
		t.Parallel()
		got := Glyph(IconModel, true)
		if got != Icons[IconModel].Glyph {
			t.Errorf("Glyph(model, true) = %q, want %q", got, Icons[IconModel].Glyph)
		}
	})

	t.Run("nerd font disabled returns fallback", func(t *testing.T) {
		t.Parallel()
		got := Glyph(IconModel, false)
		if got != Icons[IconModel].Fallback {
			t.Errorf("Glyph(model, false) = %q, want %q", got, Icons[IconModel].Fallback)
		}
	})

	t.Run("unknown key returns empty string", func(t *testing.T) {
		t.Parallel()
		if got := Glyph("not-a-real-icon", true); got != "" {
			t.Errorf("Glyph(unknown, true) = %q, want empty", got)
		}
	})
}

func TestIconsHaveNonEmptyFallbacks(t *testing.T) {
	t.Parallel()
	for key, icon := range Icons {
		if icon.Fallback == "" {
			t.Errorf("icon %q has an empty Fallback", key)
		}
	}
}
