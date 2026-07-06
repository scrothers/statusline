package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// thinkingSegment renders an icon when extended thinking is enabled for the
// session. It's a fact about the session's configuration, not a state, so
// it uses the identity accent; it's simply omitted when thinking is off or
// absent, matching the "skip the uninteresting default" pattern used by
// outputStyleSegment.
type thinkingSegment struct{}

func (thinkingSegment) ID() string { return "thinking" }

func (thinkingSegment) Priority() int { return 25 }

func (thinkingSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Thinking == nil || !rc.Payload.Thinking.Enabled {
		return nil, false
	}
	icon := theme.Glyph(theme.IconThinking, rc.Config.NerdFontEnabled())
	return []style.Chunk{{Text: icon, FG: rc.Theme.IdentityAccent}}, true
}
