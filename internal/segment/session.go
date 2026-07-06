package segment

import (
	"fmt"

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
		{Text: icon, FG: rc.Theme.Muted},
		{Text: " " + rc.Payload.SessionName, FG: rc.Theme.TextPrimary},
	}, true
}

// linesChangedSegment renders the session's total added/removed line
// counts. Omitted when there's no cost data yet, or nothing has changed.
type linesChangedSegment struct{}

func (linesChangedSegment) ID() string { return "lines_changed" }

func (linesChangedSegment) Priority() int { return 55 }

func (linesChangedSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cost := rc.Payload.Cost
	if cost == nil || (cost.TotalLinesAdded == 0 && cost.TotalLinesRemoved == 0) {
		return nil, false
	}

	var chunks []style.Chunk
	if cost.TotalLinesAdded > 0 {
		chunks = append(chunks, style.Chunk{Text: fmt.Sprintf("+%d", cost.TotalLinesAdded), FG: rc.Theme.Success})
	}
	if cost.TotalLinesRemoved > 0 {
		text := fmt.Sprintf("-%d", cost.TotalLinesRemoved)
		if len(chunks) > 0 {
			text = " " + text
		}
		chunks = append(chunks, style.Chunk{Text: text, FG: rc.Theme.Danger})
	}
	return chunks, true
}
