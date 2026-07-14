package segment

import (
	"github.com/scrothers/statusline/internal/modelid"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// providerSegment renders a small badge for the gateway a model id was
// routed through: AWS Bedrock, Google Vertex, Microsoft Foundry/Azure, a
// generic third-party router (OpenRouter and similar), or a corporate/
// self-hosted gateway. Only the glyph varies by provider — the color stays
// the theme's flat IdentityAccent, matching how repoSegment colors its
// host-branded icon — so this reads as a small supplementary badge, not a
// second accent system competing with the model segment's per-family
// colors.
//
// rc.Config.Provider is the single source of truth here: it's either an
// explicit user override (config.toml `provider = "..."`) or, more often,
// already populated by config.DetectProviderFromEnv at startup from Claude
// Code's own routing environment variables — the only reliable way to badge
// Azure/Foundry or a bare corporate-relayed id, neither of which carries any
// distinguishing shape in the id itself. modelid.DetectProvider (id-shape
// heuristics) is only the last-resort fallback when neither of those set
// anything. A plain first-party id with no gateway involved has nothing to
// detect at any tier, so the segment is omitted in the common case.
type providerSegment struct{}

func (providerSegment) ID() string { return "provider" }

// Priority is the lowest of any segment: the provider badge is the most
// supplementary signal on its line and should be the first thing dropped
// under width pressure.
func (providerSegment) Priority() int { return 15 }

func (providerSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Model == nil {
		return nil, false
	}

	provider, ok := modelid.ParseProvider(rc.Config.Provider)
	if !ok {
		provider, ok = modelid.DetectProvider(rc.Payload.Model.ID)
	}
	if !ok {
		return nil, false
	}

	iconKey := providerIconKey(provider)
	icon := theme.Glyph(iconKey, rc.Config.NerdFontEnabled())
	if icon == "" {
		return nil, false
	}
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
	}, true
}

// providerIconKey maps a detected/configured provider to its icon key.
func providerIconKey(p modelid.Provider) string {
	switch p {
	case modelid.ProviderAWS:
		return theme.IconProviderAWS
	case modelid.ProviderGCP:
		return theme.IconProviderGCP
	case modelid.ProviderAzure:
		return theme.IconProviderAzure
	case modelid.ProviderRouter:
		return theme.IconProviderRouter
	case modelid.ProviderGateway:
		return theme.IconProviderGateway
	default:
		return ""
	}
}
