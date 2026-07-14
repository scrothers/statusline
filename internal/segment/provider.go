package segment

import (
	"github.com/scrothers/statusline/internal/modelid"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// providerSegment renders a small badge for the gateway a model id was
// routed through: AWS Bedrock, Google Vertex, Microsoft Foundry/Azure, or a
// generic third-party router (OpenRouter and similar). Only the glyph
// varies by provider — the color stays the theme's flat IdentityAccent,
// matching how repoSegment colors its host-branded icon — so this reads as
// a small supplementary badge, not a second accent system competing with
// the model segment's per-family colors.
//
// A plain first-party-shaped id has no detectable provider shape, so the
// segment is omitted in the common case. Microsoft Foundry in particular
// (and especially Azure AI Foundry deployment names) can't be detected from
// id alone — it can look identical to a first-party id — so the only way
// that badge ever appears is the config.Provider override.
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
	default:
		return ""
	}
}
