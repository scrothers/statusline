package config

import "testing"

// allProviderEnvVars lists every variable DetectProviderFromEnv reads, so
// tests can force a clean baseline regardless of what's actually set in the
// host's ambient environment (this binary always runs on a real operator
// machine, not an idealized CI sandbox).
var allProviderEnvVars = []string{
	"CLAUDE_CODE_USE_BEDROCK", "ANTHROPIC_BEDROCK_BASE_URL", "ANTHROPIC_BEDROCK_MANTLE_BASE_URL",
	"AWS_BEARER_TOKEN_BEDROCK", "ANTHROPIC_AWS_BASE_URL", "ANTHROPIC_AWS_API_KEY", "ANTHROPIC_AWS_WORKSPACE_ID",
	"CLAUDE_CODE_USE_VERTEX", "ANTHROPIC_VERTEX_BASE_URL", "ANTHROPIC_VERTEX_PROJECT_ID",
	"ANTHROPIC_FOUNDRY_BASE_URL", "ANTHROPIC_FOUNDRY_RESOURCE", "ANTHROPIC_FOUNDRY_API_KEY", "ANTHROPIC_FOUNDRY_AUTH_TOKEN",
	"ANTHROPIC_BASE_URL",
}

// clearProviderEnv forces every variable DetectProviderFromEnv reads to
// empty for the duration of t, regardless of the host's ambient environment.
// t.Setenv auto-restores on cleanup, so this can't leak between tests.
func clearProviderEnv(t *testing.T) {
	t.Helper()
	for _, name := range allProviderEnvVars {
		t.Setenv(name, "")
	}
}

func TestDetectProviderFromEnv(t *testing.T) {
	// Not t.Parallel(): t.Setenv forbids it.

	t.Run("nothing set", func(t *testing.T) {
		clearProviderEnv(t)
		if _, ok := DetectProviderFromEnv(); ok {
			t.Error("DetectProviderFromEnv() ok = true, want false with nothing set")
		}
	})

	tests := []struct {
		name string
		env  string
		want string
	}{
		{"bedrock flag", "CLAUDE_CODE_USE_BEDROCK", "aws"},
		{"bedrock base url", "ANTHROPIC_BEDROCK_BASE_URL", "aws"},
		{"bedrock mantle base url", "ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "aws"},
		{"bedrock bearer token", "AWS_BEARER_TOKEN_BEDROCK", "aws"},
		{"claude platform on aws base url", "ANTHROPIC_AWS_BASE_URL", "aws"},
		{"claude platform on aws api key", "ANTHROPIC_AWS_API_KEY", "aws"},
		{"claude platform on aws workspace id", "ANTHROPIC_AWS_WORKSPACE_ID", "aws"},
		{"vertex flag", "CLAUDE_CODE_USE_VERTEX", "gcp"},
		{"vertex base url", "ANTHROPIC_VERTEX_BASE_URL", "gcp"},
		{"vertex project id", "ANTHROPIC_VERTEX_PROJECT_ID", "gcp"},
		{"foundry base url", "ANTHROPIC_FOUNDRY_BASE_URL", "azure"},
		{"foundry resource", "ANTHROPIC_FOUNDRY_RESOURCE", "azure"},
		{"foundry api key", "ANTHROPIC_FOUNDRY_API_KEY", "azure"},
		{"foundry auth token", "ANTHROPIC_FOUNDRY_AUTH_TOKEN", "azure"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearProviderEnv(t)
			t.Setenv(tt.env, "some-value")
			got, ok := DetectProviderFromEnv()
			if !ok {
				t.Fatalf("DetectProviderFromEnv() ok = false, want true with %s set", tt.env)
			}
			if got != tt.want {
				t.Errorf("DetectProviderFromEnv() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("base url pointing at the default host is not a gateway", func(t *testing.T) {
		clearProviderEnv(t)
		t.Setenv("ANTHROPIC_BASE_URL", "https://api.anthropic.com")
		if _, ok := DetectProviderFromEnv(); ok {
			t.Error("DetectProviderFromEnv() ok = true, want false for the default host")
		}
	})

	t.Run("base url pointing elsewhere is a generic gateway", func(t *testing.T) {
		clearProviderEnv(t)
		t.Setenv("ANTHROPIC_BASE_URL", "https://llm-proxy.internal.example.com")
		got, ok := DetectProviderFromEnv()
		if !ok {
			t.Fatal("DetectProviderFromEnv() ok = false, want true for a custom base URL")
		}
		if got != "gateway" {
			t.Errorf("DetectProviderFromEnv() = %q, want %q", got, "gateway")
		}
	})

	t.Run("cloud-specific signal wins over the generic base url", func(t *testing.T) {
		clearProviderEnv(t)
		t.Setenv("CLAUDE_CODE_USE_BEDROCK", "1")
		t.Setenv("ANTHROPIC_BASE_URL", "https://llm-proxy.internal.example.com")
		got, ok := DetectProviderFromEnv()
		if !ok {
			t.Fatal("DetectProviderFromEnv() ok = false, want true")
		}
		if got != "aws" {
			t.Errorf("DetectProviderFromEnv() = %q, want %q (cloud-specific signal should win)", got, "aws")
		}
	})
}
