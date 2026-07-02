package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// rateLimitGaugeWidth is narrower than the context window's hero gauge
// (contextWindowGaugeWidth) to conserve line-3 width for two of them plus
// cost, duration, and bonus badges.
const rateLimitGaugeWidth = 6

// rateLimitWindow distinguishes the two independently-absent rate-limit
// windows Claude Code reports.
type rateLimitWindow int

const (
	windowFiveHour rateLimitWindow = iota
	windowSevenDay
)

// rateLimitSegment renders one Claude subscription rate-limit gauge.
// FiveHour and SevenDay are registered as separate Segment instances
// (ratelimit_5h / ratelimit_7d) since either window may be independently
// absent from the payload.
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
	var iconKey string
	if s.window == windowFiveHour {
		win = rc.Payload.RateLimits.FiveHour
		iconKey = theme.IconRateLimitFiveHour
	} else {
		win = rc.Payload.RateLimits.SevenDay
		iconKey = theme.IconRateLimitWeek
	}
	if win == nil {
		return nil, false
	}

	bg := rc.Theme.Line3Bg
	color := gaugeColor(rc.Theme, win.UsedPercentage)
	bar := style.BlockBar(win.UsedPercentage, rateLimitGaugeWidth)
	icon := theme.Glyph(iconKey, rc.Config.NerdFontEnabled())

	return []style.Chunk{
		{Text: " " + icon + bar, FG: color, BG: bg},
		{Text: fmt.Sprintf(" %.0f%% ", win.UsedPercentage), FG: color, BG: bg},
	}, true
}
