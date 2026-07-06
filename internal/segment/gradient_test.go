package segment

import (
	"testing"

	"github.com/scrothers/statusline/internal/style"
)

func TestAggregateGradientColor(t *testing.T) {
	t.Parallel()
	th := testTheme(t)

	tests := []struct {
		name string
		pct  float64
		want style.Color
	}{
		{name: "0% is pure success", pct: 0, want: th.Success},
		{name: "50% is pure warning", pct: 50, want: th.Warning},
		{name: "100% is pure danger", pct: 100, want: th.Danger},
		{name: "negative clamps to success", pct: -10, want: th.Success},
		{name: "over 100 clamps to danger", pct: 150, want: th.Danger},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := aggregateGradientColor(&th, tt.pct); got != tt.want {
				t.Errorf("aggregateGradientColor(%v) = %+v, want %+v", tt.pct, got, tt.want)
			}
		})
	}

	t.Run("is a smooth function, not discrete bands", func(t *testing.T) {
		t.Parallel()
		// Two nearby percentages straddling the old 50% band boundary must
		// produce two DIFFERENT colors, not the same "warning" flat color a
		// discrete-band implementation would give both.
		a := aggregateGradientColor(&th, 49)
		b := aggregateGradientColor(&th, 51)
		if a == b {
			t.Errorf("aggregateGradientColor(49) == aggregateGradientColor(51) == %+v, want a smooth gradient", a)
		}
	})
}

func TestPositionGradientColor(t *testing.T) {
	t.Parallel()
	th := testTheme(t)

	tests := []struct {
		name  string
		index int
		width int
		want  style.Color
	}{
		{name: "leftmost cell is pure success", index: 0, width: 10, want: th.Success},
		{name: "rightmost cell is pure danger", index: 9, width: 10, want: th.Danger},
		{name: "midpoint cell is pure warning", index: 5, width: 11, want: th.Warning},
		{name: "single-cell bar is success", index: 0, width: 1, want: th.Success},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := positionGradientColor(&th, tt.index, tt.width); got != tt.want {
				t.Errorf("positionGradientColor(%d, %d) = %+v, want %+v", tt.index, tt.width, got, tt.want)
			}
		})
	}
}

func TestGradientBarCellChunks(t *testing.T) {
	t.Parallel()
	th := testTheme(t)

	t.Run("one chunk per filled cell, colored by position", func(t *testing.T) {
		t.Parallel()
		filled, _ := style.BlockBarParts(80, 10) // 8 full cells at 80% of width 10
		chunks := gradientBarCellChunks(&th, filled, 10)

		wantLen := len([]rune(filled))
		if len(chunks) != wantLen {
			t.Fatalf("len(chunks) = %d, want %d (one per filled rune)", len(chunks), wantLen)
		}
		for i, c := range chunks {
			want := positionGradientColor(&th, i, 10)
			if c.FG != want {
				t.Errorf("chunks[%d].FG = %+v, want %+v", i, c.FG, want)
			}
		}
	})

	t.Run("empty filled string produces no chunks", func(t *testing.T) {
		t.Parallel()
		if chunks := gradientBarCellChunks(&th, "", 10); len(chunks) != 0 {
			t.Errorf("gradientBarCellChunks(\"\", 10) = %v, want empty", chunks)
		}
	})
}
