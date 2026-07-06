package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// effortSegment renders the current reasoning effort level as an icon whose
// shape escalates with intensity — gauge glyphs from empty to full for
// low/medium/high/xhigh, then fire glyphs for max and ultra — colored along
// a fixed green-to-red-to-purple scale. Unlike most identity facts in this
// tool, effort is deliberately colored by intensity rather than the theme's
// identity accent, since "how hot is this" is the whole point of showing it.
// The scale is fixed rather than theme-derived so "purple at the top" reads
// the same regardless of which theme is active.
type effortSegment struct{}

func (effortSegment) ID() string { return "effort" }

func (effortSegment) Priority() int { return 40 }

// effortLevel pairs one Payload.Effort.Level value with its icon and its
// fixed spot on the green→red→purple intensity scale.
type effortLevel struct {
	iconKey string
	color   style.Color
}

// effortLevels covers every level level.Level currently documents
// (low/medium/high/xhigh/max) plus "ultra" for forward compatibility with a
// possible future tier beyond max.
var effortLevels = map[string]effortLevel{
	"low":    {theme.IconEffortLow, style.RGB(0x2E, 0xCC, 0x71)},    // emerald
	"medium": {theme.IconEffortMedium, style.RGB(0x8B, 0xC3, 0x4A)}, // light green
	"high":   {theme.IconEffortHigh, style.RGB(0xF1, 0xC4, 0x0F)},   // sun flower
	"xhigh":  {theme.IconEffortXHigh, style.RGB(0xE6, 0x7E, 0x22)},  // carrot
	"max":    {theme.IconEffortMax, style.RGB(0xE7, 0x4C, 0x3C)},    // alizarin
	"ultra":  {theme.IconEffortUltra, style.RGB(0x9B, 0x59, 0xB6)},  // amethyst
}

func (effortSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	level := rc.Payload.Effort
	if level == nil || level.Level == "" {
		return nil, false
	}

	entry, ok := effortLevels[level.Level]
	if !ok {
		// Unrecognized level (schema drift): still show something rather
		// than silently dropping it, in the identity accent since there's
		// no basis for picking a spot on the intensity scale.
		return []style.Chunk{{Text: level.Level, FG: rc.Theme.IdentityAccent}}, true
	}

	icon := theme.Glyph(entry.iconKey, rc.Config.NerdFontEnabled())
	return []style.Chunk{{Text: icon + " " + level.Level, FG: entry.color}}, true
}
