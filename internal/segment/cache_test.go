package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestFormatTokenCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		n    int
		want string
	}{
		{n: 0, want: "0"},
		{n: 999, want: "999"},
		{n: 1000, want: "1.0k"},
		{n: 12_300, want: "12.3k"},
		{n: 1_000_000, want: "1.0M"},
		{n: 2_500_000, want: "2.5M"},
	}
	for _, tt := range tests {
		if got := formatTokenCount(tt.n); got != tt.want {
			t.Errorf("formatTokenCount(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestCacheSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent context window is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (cacheSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil ContextWindow")
		}
	})

	t.Run("absent current usage is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{}}, nil)
		if _, ok := (cacheSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil CurrentUsage")
		}
	})

	t.Run("all-zero usage is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{}}}, nil)
		if _, ok := (cacheSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when total tokens is zero")
		}
	})

	t.Run("renders hit percentage and raw read count", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{
				CurrentUsage: &input.Usage{InputTokens: 1000, CacheCreationInputTokens: 1000, CacheReadInputTokens: 8000},
			},
		}, nil)
		chunks, ok := (cacheSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "80%") {
			t.Errorf("rendered text = %q, want it to contain 80%%", text)
		}
		if !strings.Contains(text, "8.0k") {
			t.Errorf("rendered text = %q, want it to contain 8.0k", text)
		}
	})
}
