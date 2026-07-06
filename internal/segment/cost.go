package segment

import (
	"fmt"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// costSegment renders total session cost in USD. The dollar icon is colored
// (money-green); the amount itself renders in the theme's default text
// color like ordinary prose, not a semantic accent.
type costSegment struct{}

func (costSegment) ID() string { return "cost" }

func (costSegment) Priority() int { return 70 }

func (costSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cost := rc.Payload.Cost
	if cost == nil {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	costText := fmt.Sprintf("%.2f", cost.TotalCostUSD)

	return []style.Chunk{
		{Text: theme.Glyph(theme.IconCost, nerd), FG: rc.Theme.Success},
		{Text: costText, FG: rc.Theme.TextPrimary},
	}, true
}

// durationSegment renders clock-style elapsed session duration. The clock
// icon is colored; the duration itself renders in the theme's default text
// color like ordinary prose, not a semantic accent.
type durationSegment struct{}

func (durationSegment) ID() string { return "duration" }

func (durationSegment) Priority() int { return 65 }

func (durationSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	cost := rc.Payload.Cost
	if cost == nil {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()
	durText := formatClock(cost.TotalDurationMS)

	return []style.Chunk{
		{Text: theme.Glyph(theme.IconDuration, nerd), FG: rc.Theme.Info},
		{Text: " " + durText, FG: rc.Theme.TextPrimary},
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
