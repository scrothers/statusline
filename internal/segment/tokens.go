package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// tokenCountsSegment renders a token usage breakdown behind a ticket icon:
// input/output are session-cumulative totals (how much conversation there
// is so far, the same time scope as lines_changed's cumulative add/remove
// counts), while cache-creation/cache-read are from the most recent API
// response (there's no cumulative field for either in the schema, and "how
// well is caching working right now" is inherently a per-turn question
// anyway — the same time scope the sibling cache segment already uses).
// Each category's icon carries its own color; the counts themselves are
// plain numbers in the theme's secondary text color, with no per-category
// ASCII sign. Omitted when there's no usage data at all.
type tokenCountsSegment struct{}

func (tokenCountsSegment) ID() string { return "token_counts" }

func (tokenCountsSegment) Priority() int { return 50 }

func (tokenCountsSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cw := rc.Payload.ContextWindow
	if cw == nil {
		return nil, false
	}

	var cacheCreate, cacheRead int
	if cw.CurrentUsage != nil {
		cacheCreate = cw.CurrentUsage.CacheCreationInputTokens
		cacheRead = cw.CurrentUsage.CacheReadInputTokens
	}
	if cw.TotalInputTokens == 0 && cw.TotalOutputTokens == 0 && cacheCreate == 0 && cacheRead == 0 {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	chunks := []style.Chunk{
		{Text: theme.Glyph(theme.IconTokensTicket, nerd) + "  ", FG: rc.Theme.Warning},
	}

	first := true
	for _, part := range []struct {
		iconKey string
		color   style.Color
		count   int
	}{
		{theme.IconTokensInput, rc.Theme.Success, cw.TotalInputTokens},
		{theme.IconTokensOutput, rc.Theme.Danger, cw.TotalOutputTokens},
		{theme.IconTokensCacheCreate, rc.Theme.Info, cacheCreate},
		{theme.IconCache, rc.Theme.Info, cacheRead},
	} {
		if part.count == 0 {
			continue
		}
		icon := theme.Glyph(part.iconKey, nerd)
		if !first {
			icon = " " + icon
		}
		chunks = append(chunks,
			style.Chunk{Text: icon, FG: part.color},
			style.Chunk{Text: " " + formatTokenCount(part.count), FG: rc.Theme.TextSecondary},
		)
		first = false
	}
	return chunks, true
}
