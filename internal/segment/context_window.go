package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// Bounds for the hero gauge's cell width, which scales with the detected
// terminal width (RenderContext.Columns) instead of a fixed size — rate
// limit gauges stay at a fixed, narrower width (see ratelimit.go).
const (
	contextWindowGaugeMinWidth     = 8
	contextWindowGaugeMaxWidth     = 24
	contextWindowGaugeDefaultWidth = 10 // used when Columns is unknown (<= 0)
	contextWindowGaugeDivisor      = 10 // columns per gauge cell before clamping
)

// contextAlertThreshold is the percentage at or above which the context
// gauge switches to its alarm treatment (alert icon, bold danger text),
// reserved for this one signal so it isn't diluted by other segments also
// using danger coloring.
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
	width := contextWindowGaugeWidth(rc.Columns)
	filled, track := style.BlockBarParts(pct, width)
	pctText := fmt.Sprintf("%.0f%%", pct)
	countsText := contextWindowCountsText(cw)

	if rc.Payload.Exceeds200k || pct >= contextAlertThreshold {
		icon := theme.Glyph(theme.IconContextAlert, nerd)
		text := icon + " ⟨" + filled + track + "⟩ " + pctText + countsText
		return []style.Chunk{{Text: text, FG: rc.Theme.Danger, Bold: true}}, true
	}

	color := contextGradientColor(rc.Theme, pct)
	icon := theme.Glyph(theme.IconContextWindow, nerd)
	chunks := []style.Chunk{
		{Text: icon, FG: color},
		{Text: " ⟨", FG: rc.Theme.Muted},
	}
	chunks = append(chunks, contextBarCellChunks(rc.Theme, filled, width)...)
	chunks = append(chunks,
		style.Chunk{Text: track, FG: rc.Theme.TrackDim},
		style.Chunk{Text: "⟩ ", FG: rc.Theme.Muted},
		style.Chunk{Text: pctText, FG: color},
	)
	if countsText != "" {
		chunks = append(chunks, style.Chunk{Text: countsText, FG: rc.Theme.Muted})
	}
	return chunks, true
}

// contextBarCellChunks paints each filled cell of the bar with a color
// fixed by its position along the bar — green at the left end sliding to
// red at the right — rather than one color for the whole filled run based
// on the overall percentage. That's the difference between a bar that
// reveals a stable on-screen gradient as it fills (every cell keeps the
// color it was first drawn with) and one where already-filled cells all
// shift together every time the percentage changes.
func contextBarCellChunks(th *theme.Theme, filled string, width int) []style.Chunk {
	runes := []rune(filled)
	chunks := make([]style.Chunk, len(runes))
	for i, r := range runes {
		chunks[i] = style.Chunk{Text: string(r), FG: positionGradientColor(th, i, width)}
	}
	return chunks
}

// contextWindowGaugeWidth scales the hero gauge to the detected terminal
// width instead of a fixed cell count, clamped so it's never so narrow the
// gauge is meaningless nor so wide it dominates the line on huge terminals.
func contextWindowGaugeWidth(columns int) int {
	if columns <= 0 {
		return contextWindowGaugeDefaultWidth
	}
	width := columns / contextWindowGaugeDivisor
	switch {
	case width < contextWindowGaugeMinWidth:
		return contextWindowGaugeMinWidth
	case width > contextWindowGaugeMaxWidth:
		return contextWindowGaugeMaxWidth
	default:
		return width
	}
}

// contextGradientColor returns a color that slides smoothly from the
// theme's Success (0%) through Warning (50%) to Danger (100%), rather than
// jumping between three flat bands the way gaugeColor does — used for the
// icon and percentage text, which reflect the overall (aggregate) usage.
func contextGradientColor(th *theme.Theme, pct float64) style.Color {
	return threeStopGradient(th, pct/100)
}

// positionGradientColor is contextGradientColor's counterpart for a single
// bar cell: t is the cell's fixed position (index/width-1) along the bar,
// not the current percentage, so cell 0 is always Success-ish and the last
// cell is always Danger-ish regardless of how much of the bar is filled.
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

// contextWindowCountsText renders " used/remaining" (token counts, not
// percentages) when the context window size is known, or "" when it isn't
// (e.g. before the first API response) — it's supplementary detail, so it
// simply doesn't appear rather than showing placeholder zeros.
func contextWindowCountsText(cw *input.ContextWindow) string {
	if cw.ContextWindowSize <= 0 {
		return ""
	}
	used := cw.TotalInputTokens
	remaining := max(cw.ContextWindowSize-used, 0)
	return " " + formatTokenCount(used) + "/" + formatTokenCount(remaining)
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
