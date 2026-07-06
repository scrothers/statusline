package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// tokenCountsSegment renders the most recent API response's token usage
// breakdown — input, output, cache-creation, and cache-read counts — behind
// a ticket icon. All four numbers come from the same response so they share
// one time scope; there's no per-category ASCII sign, just a distinct icon
// per category in a neutral color. Omitted when there's no usage data yet,
// or every category is zero.
type tokenCountsSegment struct{}

func (tokenCountsSegment) ID() string { return "token_counts" }

func (tokenCountsSegment) Priority() int { return 50 }

func (tokenCountsSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cw := rc.Payload.ContextWindow
	if cw == nil || cw.CurrentUsage == nil {
		return nil, false
	}

	u := cw.CurrentUsage
	if u.InputTokens == 0 && u.OutputTokens == 0 && u.CacheCreationInputTokens == 0 && u.CacheReadInputTokens == 0 {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	chunks := []style.Chunk{
		{Text: theme.Glyph(theme.IconTokensTicket, nerd), FG: rc.Theme.Warning},
	}

	for _, part := range []struct {
		iconKey string
		count   int
	}{
		{theme.IconTokensInput, u.InputTokens},
		{theme.IconTokensOutput, u.OutputTokens},
		{theme.IconTokensCacheCreate, u.CacheCreationInputTokens},
		{theme.IconCache, u.CacheReadInputTokens},
	} {
		if part.count == 0 {
			continue
		}
		icon := theme.Glyph(part.iconKey, nerd)
		chunks = append(chunks, style.Chunk{
			Text: " " + icon + " " + formatTokenCount(part.count),
			FG:   rc.Theme.Info,
		})
	}
	return chunks, true
}
