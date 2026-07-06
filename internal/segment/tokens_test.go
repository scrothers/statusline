package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestTokenCountsSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent context window is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (tokenCountsSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil ContextWindow")
		}
	})

	t.Run("zero input and output is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{}}, nil)
		if _, ok := (tokenCountsSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when both are zero")
		}
	})

	t.Run("renders both input and output with no ASCII sign", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{TotalInputTokens: 136_000, TotalOutputTokens: 8_200},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "136.0k") || !strings.Contains(text, "8.2k") {
			t.Errorf("rendered text = %q, want it to contain 136.0k and 8.2k", text)
		}
		if strings.ContainsAny(text, "+-") {
			t.Errorf("rendered text = %q, want no ASCII +/- (the icons alone carry that meaning)", text)
		}
	})

	t.Run("leads with the coin icon in warning color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{TotalInputTokens: 136_000, TotalOutputTokens: 8_200},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) == 0 {
			t.Fatal("Render() produced no chunks")
		}
		if chunks[0].FG != rc.Theme.Warning {
			t.Errorf("FG = %+v, want theme.Warning %+v", chunks[0].FG, rc.Theme.Warning)
		}
	})

	t.Run("renders input only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{TotalInputTokens: 500}}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "500") {
			t.Errorf("rendered text = %q, want it to contain 500", chunkText(chunks))
		}
	})

	t.Run("renders output only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{TotalOutputTokens: 250}}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "250") {
			t.Errorf("rendered text = %q, want it to contain 250", chunkText(chunks))
		}
	})
}
