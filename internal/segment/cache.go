package segment

import (
	"fmt"
	"strconv"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// cacheSegment renders prompt-cache effectiveness from the most recent API
// response: the fraction of input tokens served from cache, plus the raw
// cache-read token count. Unlike the context/rate-limit gauges, a high
// percentage here means the cache is doing its job, so the gradient runs
// the opposite direction: green at 100%, red at 0%.
type cacheSegment struct{}

func (cacheSegment) ID() string { return "cache" }

func (cacheSegment) Priority() int { return 35 }

func (cacheSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cw := rc.Payload.ContextWindow
	if cw == nil || cw.CurrentUsage == nil {
		return nil, false
	}

	u := cw.CurrentUsage
	total := u.InputTokens + u.CacheCreationInputTokens + u.CacheReadInputTokens
	if total == 0 {
		return nil, false
	}

	pct := float64(u.CacheReadInputTokens) / float64(total) * 100
	color := aggregateGradientColor(rc.Theme, 100-pct)
	icon := theme.Glyph(theme.IconCache, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: color},
		{Text: fmt.Sprintf(" %.0f%%", pct), FG: color},
		{Text: fmt.Sprintf(" (%s)", formatTokenCount(u.CacheReadInputTokens)), FG: rc.Theme.Muted},
	}, true
}

// formatTokenCount renders a token count with a k/M suffix for readability
// (e.g. 12300 -> "12.3k"), matching how token counts are conventionally
// displayed elsewhere in Claude Code's own UI.
func formatTokenCount(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	default:
		return strconv.Itoa(n)
	}
}
