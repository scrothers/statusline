package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// contextWindowGaugeWidth is the hero gauge's cell width; rate-limit gauges
// use a narrower width (see ratelimit.go).
const contextWindowGaugeWidth = 10

// contextAlertThreshold is the percentage at or above which the context
// gauge switches to its full-invert alarm treatment, reserved for this one
// signal so it isn't diluted by other segments also going solid-red.
const contextAlertThreshold = 95.0

// contextWindowSegment renders the context-window usage gauge: the
// session's "vitals" hero, so it renders even at 0% rather than being
// omitted when usage data is still absent.
type contextWindowSegment struct{}

func (contextWindowSegment) ID() string { return "context_window" }

func (contextWindowSegment) Priority() int { return 90 }

func (contextWindowSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cw := rc.Payload.ContextWindow
	if cw == nil {
		return nil, false
	}

	pct := contextWindowPercentage(cw)
	nerd := rc.Config.NerdFontEnabled()
	bg := rc.Theme.Line3Bg
	bar := style.BlockBar(pct, contextWindowGaugeWidth)
	pctText := fmt.Sprintf("%.0f%%", pct)

	if rc.Payload.Exceeds200k || pct >= contextAlertThreshold {
		icon := theme.Glyph(theme.IconContextAlert, nerd)
		text := " " + icon + " ⟨" + bar + "⟩ " + pctText + " "
		return []style.Chunk{{Text: text, FG: rc.Theme.Line3Bg, BG: rc.Theme.Danger, Bold: true}}, true
	}

	color := gaugeColor(rc.Theme, pct)
	icon := theme.Glyph(theme.IconContextWindow, nerd)
	return []style.Chunk{
		{Text: " " + icon, FG: color, BG: bg},
		{Text: " ⟨", FG: rc.Theme.Muted, BG: bg},
		{Text: bar, FG: color, BG: bg},
		{Text: "⟩ ", FG: rc.Theme.Muted, BG: bg},
		{Text: pctText + " ", FG: color, BG: bg},
	}, true
}

// contextWindowPercentage prefers the pre-calculated UsedPercentage, falling
// back to computing it from CurrentUsage using the same input-only formula
// Claude Code itself uses (input + cache-creation + cache-read tokens over
// the window size), and finally to 0 before the first API response.
func contextWindowPercentage(cw *input.ContextWindow) float64 {
	if cw.UsedPercentage != nil {
		return *cw.UsedPercentage
	}
	if cw.CurrentUsage != nil && cw.ContextWindowSize > 0 {
		used := cw.CurrentUsage.InputTokens + cw.CurrentUsage.CacheCreationInputTokens + cw.CurrentUsage.CacheReadInputTokens
		return float64(used) / float64(cw.ContextWindowSize) * 100
	}
	return 0
}
