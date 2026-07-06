package segment

import (
	"fmt"
	"strconv"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// cacheSegment renders prompt-cache effectiveness from the most recent API
// response: the fraction of input tokens served from cache, plus the raw
// cache-read token count.
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
	icon := theme.Glyph(theme.IconCache, rc.Config.NerdFontEnabled())
	text := fmt.Sprintf("%s %.0f%% (%s)", icon, pct, formatTokenCount(u.CacheReadInputTokens))
	return []style.Chunk{{Text: text, FG: rc.Theme.Info}}, true
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
