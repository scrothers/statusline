package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// modelSegment renders the current model's display name. It's an identity
// fact, not a state, so it always uses the theme's identity accent rather
// than a semantic success/warning/danger color.
type modelSegment struct{}

func (modelSegment) ID() string { return "model" }

// Priority is the highest of any segment: model is one of the two fields
// (with directory) that render always shrinks rather than fully drops.
func (modelSegment) Priority() int { return 100 }

func (modelSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Model == nil || rc.Payload.Model.DisplayName == "" {
		return nil, false
	}
	bg := rc.Theme.Line1Bg
	icon := theme.Glyph(theme.IconModel, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: " " + icon, FG: rc.Theme.IdentityAccent, BG: bg},
		{Text: " " + rc.Payload.Model.DisplayName + " ", FG: rc.Theme.TextPrimary, BG: bg},
	}, true
}
