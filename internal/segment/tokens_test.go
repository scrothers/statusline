package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
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

	t.Run("ticket icon has two trailing spaces for breathing room", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{CurrentUsage: &input.Usage{InputTokens: 500}},
		}, nil)
		chunks, ok := (tokenCountsSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) == 0 || !strings.HasSuffix(chunks[0].Text, "  ") {
			t.Errorf("ticket chunk = %q, want it to end with two spaces", chunks[0].Text)
		}
	})

	t.Run("each category icon carries its own color, counts use secondary text color", func(t *testing.T) {
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
		// ticket, [in-icon, in-count], [out-icon, out-count], [cc-icon, cc-count], [cr-icon, cr-count]
		if len(chunks) != 9 {
			t.Fatalf("Render() produced %d chunks, want 9", len(chunks))
		}
		wantColors := []struct {
			idx  int
			name string
			want style.Color
		}{
			{1, "input icon", rc.Theme.Success},
			{2, "input count", rc.Theme.TextSecondary},
			{3, "output icon", rc.Theme.Danger},
			{4, "output count", rc.Theme.TextSecondary},
			{5, "cache-creation icon", rc.Theme.Info},
			{6, "cache-creation count", rc.Theme.TextSecondary},
			{7, "cache-read icon", rc.Theme.Info},
			{8, "cache-read count", rc.Theme.TextSecondary},
		}
		for _, wc := range wantColors {
			if chunks[wc.idx].FG != wc.want {
				t.Errorf("%s FG = %+v, want %+v", wc.name, chunks[wc.idx].FG, wc.want)
			}
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
		if len(chunks) != 3 {
			t.Errorf("Render() produced %d chunks, want 3 (ticket + cache-read icon + cache-read count)", len(chunks))
		}
	})
}
