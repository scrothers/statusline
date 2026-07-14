package modelid

import "strings"

// Provider identifies which known gateway shape a raw model id was routed
// through.
type Provider string

// Recognized providers. ProviderAzure exists so it can be named in a config
// override — see the ParseProvider and DetectProvider docs — but is never
// returned by DetectProvider itself, since neither Microsoft Foundry model
// ids nor (especially) Azure AI Foundry deployment names carry a reliable
// distinguishing shape; a Foundry id can be indistinguishable from a
// first-party one.
const (
	ProviderAWS    Provider = "aws"
	ProviderGCP    Provider = "gcp"
	ProviderAzure  Provider = "azure"
	ProviderRouter Provider = "router"
)

// providerNames maps a config-facing provider name to its Provider value,
// used by ParseProvider.
var providerNames = map[string]Provider{
	string(ProviderAWS):    ProviderAWS,
	string(ProviderGCP):    ProviderGCP,
	string(ProviderAzure):  ProviderAzure,
	string(ProviderRouter): ProviderRouter,
}

// ParseProvider validates a user-supplied provider override string (e.g.
// from config), matching case-insensitively. It reports ok=false for an
// empty or unrecognized value, in which case the caller should fall back to
// DetectProvider.
func ParseProvider(s string) (Provider, bool) {
	p, ok := providerNames[strings.ToLower(strings.TrimSpace(s))]
	return p, ok
}

// DetectProvider inspects the raw, unmodified id — not Decode's
// progressively-stripped working string, since Decode's steps are
// destructive and order-dependent and don't expose which shape matched —
// for a known gateway signature. It reports ok=false for a bare
// first-party-shaped id, or anything else with no recognizable provider
// framing; callers should treat that as "nothing to badge," not an error.
func DetectProvider(id string) (Provider, bool) {
	s := strings.TrimSpace(id)
	if s == "" {
		return "", false
	}
	lower := strings.ToLower(s)

	// AWS Bedrock: a full ARN, or the "anthropic." provider prefix (with or
	// without a leading region, e.g. "us.anthropic.").
	if strings.Contains(lower, "arn:aws:bedrock") || reAnthropicPrefix.MatchString(s) {
		return ProviderAWS, true
	}

	// Google Vertex AI: a full resource path, or a dated-snapshot suffix.
	if strings.Contains(lower, "/publishers/") && strings.Contains(lower, "/models/") {
		return ProviderGCP, true
	}
	if i := strings.Index(s, "@"); i >= 0 && reDateToken.MatchString(s[i+1:]) {
		return ProviderGCP, true
	}

	// A generic third-party router/aggregator (OpenRouter and similar):
	// a "vendor/model" slash prefix that isn't the Vertex resource path
	// already handled above.
	if strings.Contains(s, "/") {
		return ProviderRouter, true
	}

	return "", false
}
