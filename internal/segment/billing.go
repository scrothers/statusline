package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// billingSegment renders a small badge when the session is being metered
// per-token through the API rather than drawn against a Claude subscription
// plan. input.RateLimits is populated only for Claude.ai subscribers, and
// only after their first API response in the session (see
// internal/input.RateLimits's doc comment) — so gating on Cost also being
// present rules out the early-session window where a subscriber's
// RateLimits simply hasn't arrived yet, rather than misreading it as API
// billing.
type billingSegment struct{}

func (billingSegment) ID() string { return "billing" }

// Priority sits below providerSegment's: once a session is known to be
// API-billed, which gateway it's routed through is the more useful of the
// two signals, so the billing badge is the first thing dropped under width
// pressure.
func (billingSegment) Priority() int { return 10 }

func (billingSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Cost == nil || rc.Payload.RateLimits != nil {
		return nil, false
	}

	icon := theme.Glyph(theme.IconBillingAPI, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
	}, true
}
