package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// thinkingSegment renders a lit yellow bulb when extended thinking is
// enabled, or a greyed-out bulb when it's explicitly off — unlike most
// segments, "off" is still informative here, so it's shown rather than
// omitted. Only a genuinely absent Thinking field (older clients/models
// that don't report it at all) omits the segment entirely.
type thinkingSegment struct{}

func (thinkingSegment) ID() string { return "thinking" }

func (thinkingSegment) Priority() int { return 25 }

func (thinkingSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	thinking := rc.Payload.Thinking
	if thinking == nil {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	if thinking.Enabled {
		icon := theme.Glyph(theme.IconThinkingOn, nerd)
		return []style.Chunk{{Text: icon, FG: rc.Theme.Warning}}, true
	}

	icon := theme.Glyph(theme.IconThinkingOff, nerd)
	return []style.Chunk{{Text: icon, FG: rc.Theme.Muted}}, true
}
