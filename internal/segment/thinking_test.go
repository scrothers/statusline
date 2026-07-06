package segment

import (
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestThinkingSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent thinking is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (thinkingSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Thinking")
		}
	})

	t.Run("disabled thinking is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: false}}, nil)
		if _, ok := (thinkingSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when Enabled is false")
		}
	})

	t.Run("enabled thinking renders", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: true}}, nil)
		chunks, ok := (thinkingSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) == 0 || chunkText(chunks) == "" {
			t.Error("Render() produced no content")
		}
	})
}
