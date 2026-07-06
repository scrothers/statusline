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

	t.Run("absent current usage is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{}}, nil)
		if _, ok := (tokenCountsSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil CurrentUsage")
		}
	})

	t.Run("all-zero usage is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{}}}, nil)
		if _, ok := (tokenCountsSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when every category is zero")
		}
	})

	t.Run("renders all four categories with no ASCII sign", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{
				InputTokens:              20_000,
				OutputTokens:             8_200,
				CacheCreationInputTokens: 4_000,
				CacheReadInputTokens:     108_000,
			}},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		for _, want := range []string{"20.0k", "8.2k", "4.0k", "108.0k"} {
			if !strings.Contains(text, want) {
				t.Errorf("rendered text = %q, want it to contain %q", text, want)
			}
		}
		if strings.ContainsAny(text, "+-") {
			t.Errorf("rendered text = %q, want no ASCII +/-", text)
		}
	})

	t.Run("leads with the ticket icon in warning color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{InputTokens: 500}},
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

	t.Run("token value chunks use info color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{InputTokens: 500}},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) < 2 {
			t.Fatal("Render() produced fewer than 2 chunks")
		}
		if chunks[1].FG != rc.Theme.Info {
			t.Errorf("FG = %+v, want theme.Info %+v", chunks[1].FG, rc.Theme.Info)
		}
	})

	t.Run("renders only the categories present", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{CacheReadInputTokens: 250}},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "250") {
			t.Errorf("rendered text = %q, want it to contain 250", text)
		}
		if len(chunks) != 2 {
			t.Errorf("Render() produced %d chunks, want 2 (ticket + cache-read only)", len(chunks))
		}
	})
}
