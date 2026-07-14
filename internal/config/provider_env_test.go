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

// TestDetectProviderFromEnv_baseURL is the exhaustive matrix for the
// ANTHROPIC_BASE_URL branch: default host, known gateway products, the
// generic fallback, and the malformed/adversarial-input cases a naive
// substring match would get wrong.
func TestDetectProviderFromEnv_baseURL(t *testing.T) {
	// Not t.Parallel(): t.Setenv forbids it.

	tests := []struct {
		name    string
		baseURL string
		want    string
		wantOK  bool
	}{
		{"default host exact", "https://api.anthropic.com", "", false},
		{"default host with path", "https://api.anthropic.com/v1/messages", "", false},
		{"default host different scheme case", "HTTPS://API.ANTHROPIC.COM", "", false},

		{"cloudflare ai gateway", "https://gateway.ai.cloudflare.com/v1/acct/gw/anthropic", "cloudflare", true},
		{"cloudflare bare host", "https://cloudflare.com", "cloudflare", true},
		{"cloudflare mixed case", "https://Gateway.AI.Cloudflare.COM/v1", "cloudflare", true},
		{"cloudflare with port", "https://gateway.ai.cloudflare.com:443/v1", "cloudflare", true},

		{"digitalocean inference endpoint", "https://inference.do-ai.run", "digitalocean", true},
		{"digitalocean app platform", "https://my-litellm-proxy.ondigitalocean.app", "digitalocean", true},
		{"digitalocean bare do-ai.run", "https://do-ai.run/v1", "digitalocean", true},

		{"generic corporate proxy", "https://llm-proxy.internal.mycorp.com", "gateway", true},
		{"generic proxy, bare host no scheme", "llm-gateway.mycorp.internal", "gateway", true},
		{"generic proxy with port and path", "https://10.0.0.5:8443/anthropic-proxy", "gateway", true},

		// Adversarial/false-positive guards: a naive substring match on the
		// raw string would wrongly classify all three of these.
		{"lookalike domain is not cloudflare", "https://notcloudflare.com", "gateway", true},
		{"cloudflare as a path segment is not cloudflare", "https://mycorp.com/cloudflare-proxy", "gateway", true},
		{"cloudflare as a subdomain of an unrelated domain is not cloudflare", "https://cloudflare.com.evil.example", "gateway", true},
		{"do-ai.run as a path segment is not digitalocean", "https://mycorp.com/do-ai.run", "gateway", true},
		{"api.anthropic.com in the path is not the default host", "https://evil.example/api.anthropic.com", "gateway", true},
		{"api.anthropic.com as a subdomain of an attacker domain is not the default host", "https://api.anthropic.com.evil.example", "gateway", true},

		// Malformed/edge-case input must degrade gracefully, never panic.
		{"empty", "", "", false},
		{"whitespace only", "   ", "", false},
		{"garbage no host at all", "not a url", "", false},
		{"scheme only", "https://", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearProviderEnv(t)
			t.Setenv("ANTHROPIC_BASE_URL", tt.baseURL)
			got, ok := DetectProviderFromEnv()
			if ok != tt.wantOK {
				t.Fatalf("DetectProviderFromEnv() ok = %v, want %v (base url %q)", ok, tt.wantOK, tt.baseURL)
			}
			if got != tt.want {
				t.Errorf("DetectProviderFromEnv() = %q, want %q (base url %q)", got, tt.want, tt.baseURL)
			}
		})
	}
}

func TestHostOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"full url", "https://api.anthropic.com/v1/messages", "api.anthropic.com"},
		{"uppercase host normalized to lowercase", "https://API.Anthropic.COM", "api.anthropic.com"},
		{"with port", "https://api.anthropic.com:443", "api.anthropic.com"},
		{"bare host, no scheme", "api.anthropic.com", "api.anthropic.com"},
		{"bare host with path, no scheme", "gateway.ai.cloudflare.com/v1/acct/gw", "gateway.ai.cloudflare.com"},
		{"empty", "", ""},
		{"whitespace only", "   ", ""},
		{"scheme only, no host", "https://", ""},
		{"garbage", "not a url", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := hostOf(tt.in); got != tt.want {
				t.Errorf("hostOf(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHostHasSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		host     string
		suffixes []string
		want     bool
	}{
		{"exact match", "cloudflare.com", []string{"cloudflare.com"}, true},
		{"subdomain match", "gateway.ai.cloudflare.com", []string{"cloudflare.com"}, true},
		{"deep subdomain match", "a.b.c.do-ai.run", []string{"do-ai.run"}, true},
		{"one of several suffixes matches", "inference.do-ai.run", []string{"cloudflare.com", "do-ai.run"}, true},
		{"lookalike prefix does not match", "notcloudflare.com", []string{"cloudflare.com"}, false},
		{"suffix as a subdomain of an attacker domain does not match", "cloudflare.com.evil.example", []string{"cloudflare.com"}, false},
		{"unrelated host does not match", "example.com", []string{"cloudflare.com", "do-ai.run"}, false},
		{"empty host never matches", "", []string{"cloudflare.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := hostHasSuffix(tt.host, tt.suffixes...); got != tt.want {
				t.Errorf("hostHasSuffix(%q, %v) = %v, want %v", tt.host, tt.suffixes, got, tt.want)
			}
		})
	}
}
