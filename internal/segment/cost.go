package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// costSegment renders session cost and clock-style elapsed duration.
type costSegment struct{}

func (costSegment) ID() string { return "cost" }

func (costSegment) Priority() int { return 70 }

func (costSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cost := rc.Payload.Cost
	if cost == nil {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	costText := fmt.Sprintf("$%.2f", cost.TotalCostUSD)
	durText := formatClock(cost.TotalDurationMS)

	return []style.Chunk{
		{Text: theme.Glyph(theme.IconCost, nerd) + " " + costText, FG: rc.Theme.Info},
		{Text: " " + theme.Glyph(theme.IconDuration, nerd) + " " + durText, FG: rc.Theme.Info},
	}, true
}

// formatClock renders a millisecond duration as clock-style HH:MM:SS,
// dropping the hours field (and its leading "00:") under an hour.
func formatClock(ms int64) string {
	total := ms / 1000
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
