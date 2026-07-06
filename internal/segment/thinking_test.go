package segment

import (
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestThinkingSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent thinking field is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (thinkingSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Thinking")
		}
	})

	t.Run("enabled renders the on icon in warning (yellow)", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: true}}, nil)
		chunks, ok := (thinkingSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) == 0 || chunkText(chunks) == "" {
			t.Fatal("Render() produced no content")
		}
		if chunks[0].FG != rc.Theme.Warning {
			t.Errorf("FG = %+v, want theme.Warning %+v", chunks[0].FG, rc.Theme.Warning)
		}
	})

	t.Run("explicitly disabled still renders, greyed out", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: false}}, nil)
		chunks, ok := (thinkingSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true (off is still informative, not omitted)")
		}
		if len(chunks) == 0 || chunkText(chunks) == "" {
			t.Fatal("Render() produced no content")
		}
		if chunks[0].FG != rc.Theme.Muted {
			t.Errorf("FG = %+v, want theme.Muted %+v", chunks[0].FG, rc.Theme.Muted)
		}
	})

	t.Run("on and off use different icons", func(t *testing.T) {
		t.Parallel()
		onRC := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: true}}, nil)
		offRC := newTestContext(t, &input.Payload{Thinking: &input.Thinking{Enabled: false}}, nil)

		onChunks, ok := (thinkingSegment{}).Render(onRC)
		if !ok {
			t.Fatal("Render() ok = false for enabled, want true")
		}
		offChunks, ok := (thinkingSegment{}).Render(offRC)
		if !ok {
			t.Fatal("Render() ok = false for disabled, want true")
		}
		if chunkText(onChunks) == chunkText(offChunks) {
			t.Errorf("on and off rendered identical icons: %q", chunkText(onChunks))
		}
	})
}
