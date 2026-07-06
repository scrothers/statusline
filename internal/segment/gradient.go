package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// This file holds the smooth green -> warning -> danger gradient shared by
// every gauge segment (context_window, ratelimit): a per-cell position
// gradient for the bar itself, and an aggregate-percentage gradient for the
// icon/percentage text that frames it.

// aggregateGradientColor returns a color that slides smoothly from the
// theme's Success (0%) through Warning (50%) to Danger (100%), rather than
// jumping between three flat bands — used for a gauge's icon and
// percentage text, which reflect the overall (aggregate) usage.
func aggregateGradientColor(th *theme.Theme, pct float64) style.Color {
	return threeStopGradient(th, pct/100)
}

// positionGradientColor is aggregateGradientColor's counterpart for a
// single bar cell: t is the cell's fixed position (index/(width-1)) along
// the bar, not the current percentage, so cell 0 is always Success-ish and
// the last cell is always Danger-ish regardless of how much of the bar is
// filled.
func positionGradientColor(th *theme.Theme, index, width int) style.Color {
	if width <= 1 {
		return th.Success
	}
	return threeStopGradient(th, float64(index)/float64(width-1))
}

// threeStopGradient slides from Success (t=0) through Warning (t=0.5) to
// Danger (t=1), clamping t to [0, 1].
func threeStopGradient(th *theme.Theme, t float64) style.Color {
	switch {
	case t <= 0:
		return th.Success
	case t >= 1:
		return th.Danger
	case t <= 0.5:
		return style.Lerp(th.Success, th.Warning, t/0.5)
	default:
		return style.Lerp(th.Warning, th.Danger, (t-0.5)/0.5)
	}
}

// gradientBarCellChunks paints each filled cell of a gauge bar with a color
// fixed by its position along the bar — green at the left end sliding to
// red at the right — rather than one color for the whole filled run based
// on the overall percentage. That's the difference between a bar that
// reveals a stable on-screen gradient as it fills (every cell keeps the
// color it was first drawn with) and one where already-filled cells all
// shift together every time the percentage changes.
func gradientBarCellChunks(th *theme.Theme, filled string, width int) []style.Chunk {
	runes := []rune(filled)
	chunks := make([]style.Chunk, len(runes))
	for i, r := range runes {
		chunks[i] = style.Chunk{Text: string(r), FG: positionGradientColor(th, i, width)}
	}
	return chunks
}
