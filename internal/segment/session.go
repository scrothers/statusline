package segment

import (
	"strconv"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// sessionNameSegment renders the custom session name set via --name or
// /rename. Omitted when no custom name has been set.
type sessionNameSegment struct{}

func (sessionNameSegment) ID() string { return "session_name" }

func (sessionNameSegment) Priority() int { return 45 }

func (sessionNameSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.SessionName == "" {
		return nil, false
	}
	icon := theme.Glyph(theme.IconSessionName, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
		{Text: " " + rc.Payload.SessionName, FG: rc.Theme.IdentityText},
	}, true
}

// linesChangedSegment renders the session's total added/removed line
// counts, prefixed with a pencil to mark it as an edit tally. The
// diff-added/diff-removed icons alone carry the +/- meaning, so the counts
// themselves are plain numbers with no ASCII sign. Omitted when there's no
// cost data yet, or nothing has changed.
type linesChangedSegment struct{}

func (linesChangedSegment) ID() string { return "lines_changed" }

func (linesChangedSegment) Priority() int { return 55 }

func (linesChangedSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cost := rc.Payload.Cost
	if cost == nil || (cost.TotalLinesAdded == 0 && cost.TotalLinesRemoved == 0) {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	chunks := []style.Chunk{
		{Text: theme.Glyph(theme.IconEditPencil, nerd), FG: rc.Theme.Warning},
	}
	if cost.TotalLinesAdded > 0 {
		icon := theme.Glyph(theme.IconLinesAdded, nerd)
		chunks = append(chunks, style.Chunk{
			Text: " " + icon + " " + strconv.Itoa(cost.TotalLinesAdded),
			FG:   rc.Theme.Success,
		})
	}
	if cost.TotalLinesRemoved > 0 {
		icon := theme.Glyph(theme.IconLinesRemoved, nerd)
		chunks = append(chunks, style.Chunk{
			Text: " " + icon + " " + strconv.Itoa(cost.TotalLinesRemoved),
			FG:   rc.Theme.Danger,
		})
	}
	return chunks, true
}
