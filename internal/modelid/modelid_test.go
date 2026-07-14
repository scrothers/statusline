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
