package render

import (
	"strings"

	"github.com/scrothers/statusline/internal/style"
)

// dividerGlyph is the Nerd Font thin breadcrumb-style chevron used between
// every pair of adjacent segments on a line. There is deliberately exactly
// one divider style: the statusline never paints a background, so there's
// no "pill" concept left to separate with taper caps — just plain colored
// text joined by a plain colored divider.
const dividerGlyph = "\uE0B1"

// dividerText is the full divider run between two segments: two spaces of
// padding on both sides of the glyph so segments read as clearly separate
// pieces rather than crowding the chevron.
const dividerText = "  " + dividerGlyph + "  "

// joinLine renders one line's segments into a single string: each
// segment's chunks painted in sequence, with dividerText (colored
// dividerColor) between adjacent segments. The background is always
// style.Default regardless of what a chunk's BG field holds — no segment
// may paint a background, and this is where that's enforced centrally
// rather than trusted to every segment individually. renderLine only calls
// this with a non-empty slice, having already handled the empty case
// itself.
func joinLine(segments []lineSegment, dividerColor style.Color) string {
	var b strings.Builder
	for i, seg := range segments {
		if i > 0 {
			b.WriteString(style.Paint(dividerText, dividerColor, style.Default, false))
		}
		for _, c := range seg.chunks {
			b.WriteString(style.Paint(c.Text, c.FG, style.Default, c.Bold))
		}
	}
	return b.String()
}
