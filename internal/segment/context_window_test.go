package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestContextWindowPercentage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cw   *input.ContextWindow
		want float64
	}{
		{name: "uses pre-calculated percentage", cw: &input.ContextWindow{UsedPercentage: new(float64(42))}, want: 42},
		{
			name: "falls back to current usage",
			cw: &input.ContextWindow{
				ContextWindowSize: 1000,
				CurrentUsage:      &input.Usage{InputTokens: 100, CacheCreationInputTokens: 50, CacheReadInputTokens: 50},
			},
			want: 20,
		},
		{name: "no data yet defaults to zero", cw: &input.ContextWindow{}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := contextWindowPercentage(tt.cw); got != tt.want {
				t.Errorf("contextWindowPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextWindowSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent context window is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (contextWindowSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil ContextWindow")
		}
	})

	t.Run("renders percentage", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(37))}}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "37%") {
			t.Errorf("rendered text = %q, want it to contain 37%%", chunkText(chunks))
		}
	})

	t.Run("alarm tier on exceeds_200k_tokens regardless of percentage", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(10))},
			Exceeds200k:   true,
		}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunks[0].FG != rc.Theme.Danger {
			t.Errorf("alarm tier FG = %+v, want theme.Danger %+v", chunks[0].FG, rc.Theme.Danger)
		}
		if !chunks[0].Bold {
			t.Error("alarm tier should render bold")
		}
	})

	t.Run("alarm tier at threshold percentage", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(95))}}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunks[0].FG != rc.Theme.Danger {
			t.Errorf("alarm tier FG = %+v, want theme.Danger %+v", chunks[0].FG, rc.Theme.Danger)
		}
	})

	t.Run("non-alarm tier below threshold", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(60))}}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunks[0].Bold {
			t.Error("non-alarm tier should not render bold")
		}
	})
}
