package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
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

func TestContextWindowGaugeWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		columns int
		want    int
	}{
		{name: "unknown columns uses the default", columns: 0, want: contextWindowGaugeDefaultWidth},
		{name: "negative columns uses the default", columns: -10, want: contextWindowGaugeDefaultWidth},
		{name: "narrow terminal clamps to the minimum", columns: 20, want: contextWindowGaugeMinWidth},
		{name: "huge terminal clamps to the maximum", columns: 1000, want: contextWindowGaugeMaxWidth},
		{name: "mid-size terminal scales proportionally", columns: 150, want: 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := contextWindowGaugeWidth(tt.columns); got != tt.want {
				t.Errorf("contextWindowGaugeWidth(%d) = %d, want %d", tt.columns, got, tt.want)
			}
		})
	}
}

func TestContextGradientColor(t *testing.T) {
	t.Parallel()
	th := testTheme(t)

	tests := []struct {
		name string
		pct  float64
		want style.Color
	}{
		{name: "0% is pure success", pct: 0, want: th.Success},
		{name: "50% is pure warning", pct: 50, want: th.Warning},
		{name: "100% is pure danger", pct: 100, want: th.Danger},
		{name: "negative clamps to success", pct: -10, want: th.Success},
		{name: "over 100 clamps to danger", pct: 150, want: th.Danger},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := contextGradientColor(&th, tt.pct); got != tt.want {
				t.Errorf("contextGradientColor(%v) = %+v, want %+v", tt.pct, got, tt.want)
			}
		})
	}

	t.Run("is a smooth function, not discrete bands", func(t *testing.T) {
		t.Parallel()
		// Two nearby percentages straddling the old 50% band boundary must
		// produce two DIFFERENT colors, not the same "warning" flat color a
		// discrete-band implementation would give both.
		a := contextGradientColor(&th, 49)
		b := contextGradientColor(&th, 51)
		if a == b {
			t.Errorf("contextGradientColor(49) == contextGradientColor(51) == %+v, want a smooth gradient", a)
		}
	})
}

func TestContextWindowCountsText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cw   *input.ContextWindow
		want string
	}{
		{name: "unknown window size omits counts", cw: &input.ContextWindow{}, want: ""},
		{
			name: "renders used/remaining from TotalInputTokens",
			cw:   &input.ContextWindow{ContextWindowSize: 200_000, TotalInputTokens: 50_000},
			want: " 50.0k/150.0k",
		},
		{
			name: "clamps negative remaining to zero rather than going negative",
			cw:   &input.ContextWindow{ContextWindowSize: 1000, TotalInputTokens: 5000},
			want: " 5.0k/0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := contextWindowCountsText(tt.cw); got != tt.want {
				t.Errorf("contextWindowCountsText() = %q, want %q", got, tt.want)
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

	t.Run("renders used/remaining counts when window size is known", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{
			UsedPercentage: new(float64(25)), ContextWindowSize: 200_000, TotalInputTokens: 50_000,
		}}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "50.0k/150.0k") {
			t.Errorf("rendered text = %q, want it to contain 50.0k/150.0k", chunkText(chunks))
		}
	})

	t.Run("omits counts when window size is unknown", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(25))}}, nil)
		chunks, ok := (contextWindowSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if strings.Contains(chunkText(chunks), "/") {
			t.Errorf("rendered text = %q, want no used/remaining counts", chunkText(chunks))
		}
	})

	t.Run("bar width tracks rc.Columns", func(t *testing.T) {
		t.Parallel()
		payload := &input.Payload{ContextWindow: &input.ContextWindow{UsedPercentage: new(float64(100))}}
		narrow := newTestContext(t, payload, nil)
		narrow.Columns = 40 // width -> min clamp of 8
		wide := newTestContext(t, payload, nil)
		wide.Columns = 300 // width -> max clamp of 24

		narrowChunks, ok := (contextWindowSegment{}).Render(narrow)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		wideChunks, ok := (contextWindowSegment{}).Render(wide)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunkText(narrowChunks)) >= len(chunkText(wideChunks)) {
			t.Errorf("expected a wider terminal to produce a longer bar: narrow=%q wide=%q",
				chunkText(narrowChunks), chunkText(wideChunks))
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
