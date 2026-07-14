package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

// disableNerdFont points rc.Config.NerdFont at false so rendered chunks use
// the plain-text fallback instead of a Nerd Font glyph, making assertions
// on rendered text portable regardless of the default config.
func disableNerdFont(rc *RenderContext) {
	nerdFont := false
	rc.Config.NerdFont = &nerdFont
}

func TestProviderSegment(t *testing.T) {
	t.Parallel()

	t.Run("renders a badge for a detectable bedrock id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID: "anthropic.claude-opus-4-8",
		}}, nil)
		disableNerdFont(rc)
		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "AWS") {
			t.Errorf("rendered text = %q, want it to contain the AWS fallback", chunkText(chunks))
		}
	})

	t.Run("renders a badge for a detectable vertex id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID: "claude-opus-4-5@20251101",
		}}, nil)
		disableNerdFont(rc)
		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "GCP") {
			t.Errorf("rendered text = %q, want it to contain the GCP fallback", chunkText(chunks))
		}
	})

	t.Run("renders a badge for an openrouter-style id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID: "anthropic/claude-3.5-sonnet:beta",
		}}, nil)
		disableNerdFont(rc)
		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "Router") {
			t.Errorf("rendered text = %q, want it to contain the Router fallback", chunkText(chunks))
		}
	})

	t.Run("omitted for a plain first-party id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{ID: "claude-opus-4-8"}}, nil)
		_, ok := providerSegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false for a plain first-party id")
		}
	})

	t.Run("omitted for an absent model", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		_, ok := providerSegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false for nil Model")
		}
	})

	t.Run("config override forces the gateway badge on an otherwise-plain id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{ID: "claude-4-8-opus"}}, nil)
		disableNerdFont(rc)
		rc.Config.Provider = "gateway"

		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true with a provider override set")
		}
		if !strings.Contains(chunkText(chunks), "Gateway") {
			t.Errorf("rendered text = %q, want it to contain the Gateway fallback", chunkText(chunks))
		}
	})

	t.Run("config override forces the azure badge on an otherwise-plain id", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{ID: "claude-opus-4-8"}}, nil)
		disableNerdFont(rc)
		rc.Config.Provider = "azure"

		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true with a provider override set")
		}
		if !strings.Contains(chunkText(chunks), "Azure") {
			t.Errorf("rendered text = %q, want it to contain the Azure fallback", chunkText(chunks))
		}
	})

	t.Run("config override wins over auto-detection", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID: "anthropic.claude-opus-4-8", // would auto-detect as AWS
		}}, nil)
		disableNerdFont(rc)
		rc.Config.Provider = "router"

		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "Router") {
			t.Errorf("rendered text = %q, want the override to win over auto-detection", chunkText(chunks))
		}
	})

	t.Run("unrecognized config override falls back to auto-detection", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Model: &input.Model{
			ID: "anthropic.claude-opus-4-8",
		}}, nil)
		disableNerdFont(rc)
		rc.Config.Provider = "not-a-real-provider"

		chunks, ok := providerSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true (auto-detection should still run)")
		}
		if !strings.Contains(chunkText(chunks), "AWS") {
			t.Errorf("rendered text = %q, want auto-detected AWS", chunkText(chunks))
		}
	})
}
