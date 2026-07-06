package style

import (
	"testing"
	"unicode/utf8"
)

func TestBlockBarParts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pct        float64
		width      int
		wantFilled string
		wantTrack  string
	}{
		{name: "zero percent is fully empty", pct: 0, width: 5, wantFilled: "", wantTrack: "░░░░░"},
		{name: "hundred percent is fully filled", pct: 100, width: 5, wantFilled: "█████", wantTrack: ""},
		{name: "half of ten cells", pct: 50, width: 10, wantFilled: "█████", wantTrack: "░░░░░"},
		{name: "negative clamps to zero", pct: -10, width: 3, wantFilled: "", wantTrack: "░░░"},
		{name: "over 100 clamps to full", pct: 150, width: 3, wantFilled: "███", wantTrack: ""},
		{name: "zero width is empty strings", pct: 50, width: 0, wantFilled: "", wantTrack: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filled, track := BlockBarParts(tt.pct, tt.width)
			if filled != tt.wantFilled || track != tt.wantTrack {
				t.Errorf("BlockBarParts(%v, %d) = (%q, %q), want (%q, %q)",
					tt.pct, tt.width, filled, track, tt.wantFilled, tt.wantTrack)
			}
		})
	}
}

func TestBlockBarParts_alwaysWidthCellsWideCombined(t *testing.T) {
	t.Parallel()
	for pct := 0.0; pct <= 100.0; pct += 3.7 {
		filled, track := BlockBarParts(pct, 10)
		combined := filled + track
		if n := utf8.RuneCountInString(combined); n != 10 {
			t.Errorf("BlockBarParts(%v, 10) combined has %d cells, want 10 (got %q)", pct, n, combined)
		}
	}
}
