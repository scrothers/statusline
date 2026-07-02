package render

import (
	"strings"

	"github.com/scrothers/statusline/internal/style"
)

// capGlyphs is the opening/connector-and-closing glyph pair for one
// powerline separator style. The connector and closing cap share a glyph;
// only fg/bg differ (see joinLine).
type capGlyphs struct {
	Open  string
	Close string
}

// defaultCapStyle is used when Config.Separator.Style is empty or doesn't
// match a known style.
const defaultCapStyle = "rounded"

var capStyles = map[string]capGlyphs{
	"rounded": {Open: "\uE0B6", Close: "\uE0B4"},
	"hard":    {Open: "\uE0B2", Close: "\uE0B0"},
}

func capsFor(styleName string) capGlyphs {
	if g, ok := capStyles[styleName]; ok {
		return g
	}
	return capStyles[defaultCapStyle]
}

// joinLine renders one line's segments into a single string, gluing them
// together based on each segment's background color:
//
//   - pill → pill (both bg.Valid): a connector glyph colored
//     fg=leaving segment's bg, bg=entering segment's bg. Getting this
//     backwards (fg/bg swapped) draws a flag/pennant instead of a
//     continuous ribbon.
//   - pill → badge or badge → pill: a closing/opening cap tapering to the
//     terminal's own background, plus a plain space, since a badge has no
//     background to connect to.
//   - badge → badge (neither bg.Valid): a plain " · " divider in muted.
//
// The first segment gets a leading opening cap (if it's a pill) and the
// last gets a trailing closing cap tapering to the terminal default,
// keeping the ribbon self-contained regardless of terminal background.
func joinLine(segments []lineSegment, capStyleName string, muted style.Color) string {
	if len(segments) == 0 {
		return ""
	}
	caps := capsFor(capStyleName)

	var b strings.Builder
	for i, seg := range segments {
		if i == 0 {
			if seg.bg.Valid {
				b.WriteString(style.Paint(caps.Open, seg.bg, style.Default, false))
			}
		} else {
			prev := segments[i-1]
			switch {
			case prev.bg.Valid && seg.bg.Valid:
				b.WriteString(style.Paint(caps.Close, prev.bg, seg.bg, false))
			case prev.bg.Valid && !seg.bg.Valid:
				b.WriteString(style.Paint(caps.Close, prev.bg, style.Default, false))
				b.WriteString(" ")
			case !prev.bg.Valid && seg.bg.Valid:
				b.WriteString(" ")
				b.WriteString(style.Paint(caps.Open, seg.bg, style.Default, false))
			default:
				b.WriteString(style.Paint(" · ", muted, style.Default, false))
			}
		}
		for _, c := range seg.chunks {
			b.WriteString(style.Paint(c.Text, c.FG, c.BG, c.Bold))
		}
	}

	last := segments[len(segments)-1]
	if last.bg.Valid {
		b.WriteString(style.Paint(caps.Close, last.bg, style.Default, false))
	}
	return b.String()
}
