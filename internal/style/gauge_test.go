package style

import (
	"testing"
	"unicode/utf8"
)

func TestBlockBar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		pct   float64
		width int
		want  string
	}{
		{name: "zero percent is fully empty", pct: 0, width: 5, want: "░░░░░"},
		{name: "hundred percent is fully filled", pct: 100, width: 5, want: "█████"},
		{name: "half of ten cells", pct: 50, width: 10, want: "█████░░░░░"},
		{name: "negative clamps to zero", pct: -10, width: 3, want: "░░░"},
		{name: "over 100 clamps to full", pct: 150, width: 3, want: "███"},
		{name: "zero width is empty string", pct: 50, width: 0, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := BlockBar(tt.pct, tt.width); got != tt.want {
				t.Errorf("BlockBar(%v, %d) = %q, want %q", tt.pct, tt.width, got, tt.want)
			}
		})
	}
}

func TestBlockBar_alwaysWidthCellsWide(t *testing.T) {
	t.Parallel()
	for pct := 0.0; pct <= 100.0; pct += 3.7 {
		got := BlockBar(pct, 10)
		if n := utf8.RuneCountInString(got); n != 10 {
			t.Errorf("BlockBar(%v, 10) has %d cells, want 10 (got %q)", pct, n, got)
		}
	}
}
