package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// tokenCountsSegment renders the session's total input/output token counts,
// prefixed with a coin icon and formatted the same way as
// linesChangedSegment's added/removed counters: the diff-added/diff-removed
// icons carry the +/- meaning, the numbers are plain. Omitted when there's
// no context-window data yet, or nothing has been read/written.
type tokenCountsSegment struct{}

func (tokenCountsSegment) ID() string { return "token_counts" }

func (tokenCountsSegment) Priority() int { return 50 }

func (tokenCountsSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cw := rc.Payload.ContextWindow
	if cw == nil || (cw.TotalInputTokens == 0 && cw.TotalOutputTokens == 0) {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	chunks := []style.Chunk{
		{Text: theme.Glyph(theme.IconTokensCoin, nerd), FG: rc.Theme.Warning},
	}
	if cw.TotalInputTokens > 0 {
		icon := theme.Glyph(theme.IconLinesAdded, nerd)
		chunks = append(chunks, style.Chunk{
			Text: " " + icon + " " + formatTokenCount(cw.TotalInputTokens),
			FG:   rc.Theme.Success,
		})
	}
	if cw.TotalOutputTokens > 0 {
		icon := theme.Glyph(theme.IconLinesRemoved, nerd)
		chunks = append(chunks, style.Chunk{
			Text: " " + icon + " " + formatTokenCount(cw.TotalOutputTokens),
			FG:   rc.Theme.Danger,
		})
	}
	return chunks, true
}
