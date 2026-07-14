package segment

import (
	"github.com/scrothers/statusline/internal/modelid"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// modelSegment renders the current model's display name. It's an identity
// fact, not a state, so its icon uses a per-model-family identity accent
// (Theme.IdentityColorFor) rather than a semantic success/warning/danger
// color — the label text itself stays a neutral, theme-consistent color.
type modelSegment struct{}

func (modelSegment) ID() string { return "model" }

// Priority is the highest of any segment: model is one of the two fields
// (with directory) that render always shrinks rather than fully drops.
func (modelSegment) Priority() int { return 100 }

func (modelSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Model == nil {
		return nil, false
	}
	label := modelid.Label(rc.Payload.Model.ID, rc.Payload.Model.DisplayName)
	if label == "" {
		return nil, false
	}
	family, _ := modelid.Family(rc.Payload.Model.ID)
	icon := theme.Glyph(theme.IconModel, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityColorFor(family)},
		{Text: " " + label, FG: rc.Theme.IdentityText},
	}, true
}
