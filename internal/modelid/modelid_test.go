package modelid

import "testing"

func TestDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		wantLabel string
		wantOK    bool
	}{
		{"first-party opus", "claude-opus-4-8", "Opus 4.8", true},
		{"gateway reordered opus", "claude-4-8-opus", "Opus 4.8", true},
		{"gateway reordered sonnet", "claude-5-sonnet", "Sonnet 5", true},
		{"gateway reordered haiku", "claude-4-5-haiku", "Haiku 4.5", true},
		{"1m marker dropped", "claude-opus-4-8[1m]", "Opus 4.8", true},
		{"1m marker uppercase", "claude-opus-4-8[1M]", "Opus 4.8", true},
		{"gateway reordered with 1m", "claude-4-8-opus[1m]", "Opus 4.8", true},
		{
			"bedrock cross-region dated",
			"us.anthropic.claude-3-5-sonnet-20241022-v2:0",
			"Sonnet 3.5",
			true,
		},
		{"bedrock plain", "anthropic.claude-opus-4-8", "Opus 4.8", true},
		{"vertex dated snapshot", "claude-opus-4-5@20251101", "Opus 4.5", true},
		{"openrouter dotted tagged", "anthropic/claude-3.5-sonnet:beta", "Sonnet 3.5", true},
		{"full pinned id", "claude-haiku-4-5-20251001", "Haiku 4.5", true},
		{"fable no minor version", "claude-fable-5", "Fable 5", true},
		{"mythos", "claude-mythos-5", "Mythos 5", true},
		{"legacy dated", "claude-3-opus-20240229", "Opus 3", true},
		{"legacy dotted no family word", "claude-2.1", "Claude 2.1", true},
		{"family word only, no version", "claude-opus", "Opus", true},
		{
			"bedrock full arn",
			"arn:aws:bedrock:us-east-1::foundation-model/anthropic.claude-3-5-sonnet-20241022-v2:0",
			"Sonnet 3.5",
			true,
		},
		{
			"bedrock inference-profile arn",
			"arn:aws:bedrock:us-east-1:123456789012:inference-profile/us.anthropic.claude-3-5-sonnet-20241022-v2:0",
			"Sonnet 3.5",
			true,
		},
		{
			"vertex full resource path",
			"projects/my-proj/locations/us-central1/publishers/anthropic/models/claude-opus-4-5@20251101",
			"Opus 4.5",
			true,
		},
		{"underscore separated", "claude_opus_4_8", "Opus 4.8", true},
		{"whitespace separated", "claude opus 4 8", "Opus 4.8", true},
		{"mixed case whole id", "Claude-Opus-4-8", "Opus 4.8", true},
		{"uppercase anthropic prefix", "ANTHROPIC.claude-opus-4-8", "Opus 4.8", true},
		{"legacy zero-alias opus", "claude-opus-4-0", "Opus 4", true},
		{"legacy zero-alias sonnet", "claude-sonnet-4-0", "Sonnet 4", true},
		{"filler word with real family", "claude-3-5-sonnet-beta", "Sonnet 3.5", true},
		{"mythos preview alias", "claude-mythos-preview", "Mythos", true},
		{"stray build-number token ignored", "claude-opus-4-8-001", "Opus 4.8", true},
		{"filler word alone falls back to claude", "claude-latest", "Claude", true},
		{"bare 1m marker alone falls back to claude", "claude-1m", "Claude", true},
		{"unknown future family word", "claude-atlas-6", "Atlas 6", true},
		{"empty id", "", "", false},
		{"whitespace only", "   ", "", false},
		{"non-claude id", "gpt-4-turbo", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotLabel, gotOK := Decode(tt.id)
			if gotOK != tt.wantOK {
				t.Fatalf("Decode(%q) ok = %v, want %v", tt.id, gotOK, tt.wantOK)
			}
			if gotLabel != tt.wantLabel {
				t.Errorf("Decode(%q) label = %q, want %q", tt.id, gotLabel, tt.wantLabel)
			}
		})
	}
}

func TestFamily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         string
		wantFamily string
		wantOK     bool
	}{
		{"opus", "claude-opus-4-8", "Opus", true},
		{"gateway reordered opus", "claude-4-8-opus", "Opus", true},
		{"sonnet", "claude-sonnet-4-6", "Sonnet", true},
		{"haiku", "claude-haiku-4-5", "Haiku", true},
		{"fable", "claude-fable-5", "Fable", true},
		{"mythos", "claude-mythos-5", "Mythos", true},
		{"legacy dotted no family word falls back to claude", "claude-2.1", "Claude", true},
		{"unknown future family word", "claude-atlas-6", "Atlas", true},
		{"bedrock wrapped sonnet", "anthropic.claude-3-5-sonnet-20241022", "Sonnet", true},
		{"empty id", "", "", false},
		{"non-claude id", "gpt-4-turbo", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotFamily, gotOK := Family(tt.id)
			if gotOK != tt.wantOK {
				t.Fatalf("Family(%q) ok = %v, want %v", tt.id, gotOK, tt.wantOK)
			}
			if gotFamily != tt.wantFamily {
				t.Errorf("Family(%q) = %q, want %q", tt.id, gotFamily, tt.wantFamily)
			}
		})
	}
}

func TestLabel(t *testing.T) {
	t.Parallel()

	t.Run("decodes a recognizable id", func(t *testing.T) {
		t.Parallel()
		got := Label("claude-4-8-opus", "whatever the gateway sent")
		if got != "Opus 4.8" {
			t.Errorf("Label() = %q, want %q", got, "Opus 4.8")
		}
	})

	t.Run("falls back to displayName when id doesn't decode", func(t *testing.T) {
		t.Parallel()
		got := Label("gpt-4-turbo", "Some Display Name")
		if got != "Some Display Name" {
			t.Errorf("Label() = %q, want %q", got, "Some Display Name")
		}
	})

	t.Run("falls back to displayName when id is empty", func(t *testing.T) {
		t.Parallel()
		got := Label("", "Opus")
		if got != "Opus" {
			t.Errorf("Label() = %q, want %q", got, "Opus")
		}
	})

	t.Run("empty when both id and displayName are unusable", func(t *testing.T) {
		t.Parallel()
		got := Label("", "")
		if got != "" {
			t.Errorf("Label() = %q, want empty", got)
		}
	})
}
