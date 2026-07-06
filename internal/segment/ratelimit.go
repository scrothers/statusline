package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// rateLimitGaugeWidth is narrower than the context window's hero gauge
// (which scales with terminal width) to conserve room for the rest of the
// Claude info line.
const rateLimitGaugeWidth = 6

// rateLimitWindow distinguishes the two independently-absent rate-limit
// windows Claude Code reports.
type rateLimitWindow int

const (
	windowFiveHour rateLimitWindow = iota
	windowSevenDay
)

// rateLimitSegment renders one Claude subscription rate-limit gauge, styled
// to match the context-window gauge: an icon, an explicit window label
// ("5h"/"7d") so the two gauges aren't distinguishable by icon alone, a
// bracketed bar whose filled cells follow the same fixed green-to-red
// position gradient, and the percentage. FiveHour and SevenDay are
// registered as separate Segment instances (ratelimit_5h / ratelimit_7d)
// since either window may be independently absent from the payload.
type rateLimitSegment struct {
	window rateLimitWindow
}

func newRateLimitSegment(w rateLimitWindow) rateLimitSegment {
	return rateLimitSegment{window: w}
}

func (s rateLimitSegment) ID() string {
	if s.window == windowFiveHour {
		return "ratelimit_5h"
	}
	return "ratelimit_7d"
}

func (rateLimitSegment) Priority() int { return 50 }

func (s rateLimitSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.RateLimits == nil {
		return nil, false
	}

	var win *input.RateLimitWindow
	var iconKey, label string
	if s.window == windowFiveHour {
		win = rc.Payload.RateLimits.FiveHour
		iconKey = theme.IconRateLimitFiveHour
		label = "5h"
	} else {
		win = rc.Payload.RateLimits.SevenDay
		iconKey = theme.IconRateLimitWeek
		label = "7d"
	}
	if win == nil {
		return nil, false
	}

	color := aggregateGradientColor(rc.Theme, win.UsedPercentage)
	filled, track := style.BlockBarParts(win.UsedPercentage, rateLimitGaugeWidth)
	icon := theme.Glyph(iconKey, rc.Config.NerdFontEnabled())

	chunks := make([]style.Chunk, 0, rateLimitGaugeWidth+5)
	chunks = append(chunks,
		style.Chunk{Text: icon + " " + label, FG: color},
		style.Chunk{Text: " ⟨", FG: rc.Theme.Muted},
	)
	chunks = append(chunks, gradientBarCellChunks(rc.Theme, filled, rateLimitGaugeWidth)...)
	chunks = append(chunks,
		style.Chunk{Text: track, FG: rc.Theme.TrackDim},
		style.Chunk{Text: "⟩", FG: rc.Theme.Muted},
		style.Chunk{Text: fmt.Sprintf(" %.0f%%", win.UsedPercentage), FG: color},
	)
	return chunks, true
}
