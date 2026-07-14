package config

import (
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
// proxy — a corporate/self-hosted gateway being the case none of the
// cloud-specific variables above already explain.
const defaultAnthropicHost = "api.anthropic.com"

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

	if base := os.Getenv("ANTHROPIC_BASE_URL"); base != "" &&
		!strings.Contains(strings.ToLower(base), defaultAnthropicHost) {
		return "gateway", true
	}

	return "", false
}
