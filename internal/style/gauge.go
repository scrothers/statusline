package style

import "strings"

// blockLevels are the seven sub-cell fill glyphs, from one-eighth to
// seven-eighths full; a fully-filled cell uses '█' directly.
var blockLevels = [...]rune{'▏', '▎', '▍', '▌', '▋', '▊', '▉'}

// BlockBar renders pct (clamped to [0, 100]) as a width-cell bar using an
// eight-level smooth block ramp, giving width*8 discrete fill levels instead
// of the width levels a plain filled/empty bar would give. The returned
// string carries no color; callers paint it via Chunk.FG.
func BlockBar(pct float64, width int) string {
	if width <= 0 {
		return ""
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}

	maxEighths := width * 8
	eighths := min(int(pct/100*float64(maxEighths)+0.5), maxEighths)
	full, partial := eighths/8, eighths%8

	var b strings.Builder
	for range full {
		b.WriteRune('█')
	}
	if partial > 0 && full < width {
		b.WriteRune(blockLevels[partial-1])
		full++
	}
	for range width - full {
		b.WriteRune('░')
	}
	return b.String()
}
