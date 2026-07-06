package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestFormatClock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ms   int64
		want string
	}{
		{ms: 0, want: "00:00"},
		{ms: 45_000, want: "00:45"},
		{ms: 125_000, want: "02:05"},
		{ms: 3_661_000, want: "1:01:01"},
		{ms: 3_600_000, want: "1:00:00"},
	}

	for _, tt := range tests {
		if got := formatClock(tt.ms); got != tt.want {
			t.Errorf("formatClock(%d) = %q, want %q", tt.ms, got, tt.want)
		}
	}
}

func TestCostSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent cost is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (costSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Cost")
		}
	})

	t.Run("renders cost", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalCostUSD: 1.234}}, nil)
		chunks, ok := (costSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "1.23") {
			t.Errorf("rendered text = %q, want it to contain 1.23", text)
		}
		if strings.Count(text, "$") > 0 {
			t.Errorf("rendered text = %q, want no literal $ when the nerd font dollar icon already renders one", text)
		}
		if strings.Contains(text, " 1.23") {
			t.Errorf("rendered text = %q, want no space between the dollar icon and the amount", text)
		}
	})

	t.Run("plain fallback shows exactly one dollar sign directly against the amount", func(t *testing.T) {
		t.Parallel()
		nerdFont := false
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalCostUSD: 1.234}}, nil)
		rc.Config.NerdFont = &nerdFont
		chunks, ok := (costSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if strings.Count(text, "$") != 1 {
			t.Errorf("rendered text = %q, want exactly one $", text)
		}
		if !strings.Contains(text, "$1.23") {
			t.Errorf("rendered text = %q, want $ directly against 1.23 with no space", text)
		}
	})
}

func TestDurationSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent cost is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (durationSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Cost")
		}
	})

	t.Run("renders clock-style duration", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalDurationMS: 65_000}}, nil)
		chunks, ok := (durationSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "01:05") {
			t.Errorf("rendered text = %q, want it to contain 01:05", chunkText(chunks))
		}
	})
}
