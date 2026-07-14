package config

import (
	"net/url"
	"os"
	"strings"
)

// providerEnvVars lists, for each provider, the Claude Code environment
// variables whose presence indicates traffic is routed that way — sourced
// from Claude Code's documented routing env vars (code.claude.com/docs/en/env-vars).
// Checked in this order (first match wins) so a var naming a specific cloud
// is trusted over the generic base-URL override below.
var providerEnvVars = []struct {
	provider string
	names    []string
}{
	{"aws", []string{
		"CLAUDE_CODE_USE_BEDROCK",
		"ANTHROPIC_BEDROCK_BASE_URL",
		"ANTHROPIC_BEDROCK_MANTLE_BASE_URL",
		"AWS_BEARER_TOKEN_BEDROCK",
		"ANTHROPIC_AWS_BASE_URL",
		"ANTHROPIC_AWS_API_KEY",
		"ANTHROPIC_AWS_WORKSPACE_ID",
	}},
	{"gcp", []string{
		"CLAUDE_CODE_USE_VERTEX",
		"ANTHROPIC_VERTEX_BASE_URL",
		"ANTHROPIC_VERTEX_PROJECT_ID",
	}},
	{"azure", []string{
		"ANTHROPIC_FOUNDRY_BASE_URL",
		"ANTHROPIC_FOUNDRY_RESOURCE",
		"ANTHROPIC_FOUNDRY_API_KEY",
		"ANTHROPIC_FOUNDRY_AUTH_TOKEN",
	}},
}

// defaultAnthropicHost is the first-party API host. ANTHROPIC_BASE_URL
// pointing anywhere else means traffic is routed through some kind of
// proxy or gateway.
const defaultAnthropicHost = "api.anthropic.com"

// gatewayHosts maps a specific, verified gateway-product hostname suffix to
// its provider, checked before falling back to the generic "gateway" badge.
// Sourced from each product's own documentation:
//   - Cloudflare AI Gateway: gateway.ai.cloudflare.com (developers.cloudflare.com/ai-gateway)
//   - DigitalOcean Gradient AI: inference.do-ai.run, or any app hosted on
//     DigitalOcean's App Platform (*.ondigitalocean.app) proxying it
//
// Suffixes, not exact hosts, so a subdomain still matches (e.g. a
// per-account gateway.ai.cloudflare.com path segment doesn't change the
// host, but a future product variant on a subdomain would still resolve
// correctly here).
var gatewayHosts = []struct {
	provider string
	suffixes []string
}{
	{"cloudflare", []string{"cloudflare.com"}},
	{"digitalocean", []string{"do-ai.run", "ondigitalocean.app"}},
}

// DetectProviderFromEnv inspects Claude Code's own routing environment
// variables (inherited from the parent process, since this binary always
// runs as a Claude Code subprocess) for a known gateway signal, and reports
// a provider string suitable for Config.Provider. This is a stronger signal
// than guessing from the model id: it's Claude Code's own routing decision,
// not an inference — and it's the only way to detect Azure/Foundry or a
// generic corporate gateway at all, since neither carries any distinguishing
// shape in the id itself. It reports ok=false when nothing is set, in which
// case the caller should leave Config.Provider for id-shape auto-detection.
func DetectProviderFromEnv() (string, bool) {
	for _, group := range providerEnvVars {
		for _, name := range group.names {
			if os.Getenv(name) != "" {
				return group.provider, true
			}
		}
	}

	base := os.Getenv("ANTHROPIC_BASE_URL")
	if base == "" {
		return "", false
	}
	host := hostOf(base)
	if host == "" {
		// Set but unparseable as a URL/host — can't rule out the default
		// host, so there's nothing safe to report rather than a guess.
		return "", false
	}
	if host == defaultAnthropicHost {
		return "", false
	}
	for _, gw := range gatewayHosts {
		if hostHasSuffix(host, gw.suffixes...) {
			return gw.provider, true
		}
	}
	return "gateway", true
}

// hostOf returns the lowercase hostname of rawURL, or "" if it can't be
// parsed as one — never panics on malformed input. Handles a bare
// "host[:port][/path]" value with no scheme at all (a plausible typo/paste
// error in an env var, e.g. "gateway.ai.cloudflare.com" with no leading
// "https://") by retrying with an assumed scheme, rather than silently
// treating the whole string as an unparseable relative path.
//
// The retry only fires when rawURL has no scheme — not merely an empty
// host. A string like "https://" already has a scheme and a (deliberately)
// empty host; retrying that by prepending another "https://" produces the
// nonsense "https://https://", which Go happily mis-parses a bogus non-empty
// host out of ("https"). Gating on Scheme=="" avoids that trap.
func hostOf(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return ""
	}
	if u.Scheme == "" && u.Host == "" {
		if u2, err2 := url.Parse("https://" + trimmed); err2 == nil {
			u = u2
		}
	}
	return strings.ToLower(u.Hostname())
}

// hostHasSuffix reports whether host is exactly one of suffixes, or a
// subdomain of one — a dot-boundary suffix match, not a raw substring
// match, so "notcloudflare.com" or "cloudflare.com.evil.example" never
// falsely match "cloudflare.com".
func hostHasSuffix(host string, suffixes ...string) bool {
	for _, suf := range suffixes {
		if host == suf || strings.HasSuffix(host, "."+suf) {
			return true
		}
	}
	return false
}
