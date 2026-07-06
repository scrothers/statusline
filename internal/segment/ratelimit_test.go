package segment

import (
	"strings"
	"testing"
	"time"

	"github.com/scrothers/statusline/internal/input"
)

func TestRateLimitSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent rate limits is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := newRateLimitSegment(windowFiveHour).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil RateLimits")
		}
	})

	t.Run("one window present, the other absent", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{FiveHour: &input.RateLimitWindow{UsedPercentage: 23.5}},
		}, nil)

		chunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("five_hour Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "24%") && !strings.Contains(chunkText(chunks), "23%") {
			t.Errorf("rendered text = %q, want it to contain the rounded percentage", chunkText(chunks))
		}

		if _, ok := newRateLimitSegment(windowSevenDay).Render(rc); ok {
			t.Error("seven_day Render() ok = true, want false (absent)")
		}
	})

	t.Run("IDs are distinct", func(t *testing.T) {
		t.Parallel()
		if newRateLimitSegment(windowFiveHour).ID() == newRateLimitSegment(windowSevenDay).ID() {
			t.Error("five-hour and seven-day segments share an ID")
		}
	})

	t.Run("labels distinguish the two windows", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 10},
				SevenDay: &input.RateLimitWindow{UsedPercentage: 10},
			},
		}, nil)

		fiveHourChunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("five_hour Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(fiveHourChunks), "5h") {
			t.Errorf("five-hour rendered text = %q, want it to contain the 5h label", chunkText(fiveHourChunks))
		}

		sevenDayChunks, ok := newRateLimitSegment(windowSevenDay).Render(rc)
		if !ok {
			t.Fatal("seven_day Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(sevenDayChunks), "7d") {
			t.Errorf("seven-day rendered text = %q, want it to contain the 7d label", chunkText(sevenDayChunks))
		}
	})

	t.Run("no reset time is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{FiveHour: &input.RateLimitWindow{UsedPercentage: 10}},
		}, nil)

		chunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if strings.Contains(chunkText(chunks), "resets") {
			t.Errorf("rendered text = %q, want no reset text when ResetsAt is zero", chunkText(chunks))
		}
	})

	t.Run("renders a muted reset countdown when ResetsAt is present", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		resetsAt := now.Add(5 * time.Hour)
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 10, ResetsAt: resetsAt.Unix()},
			},
		}, nil)
		rc.Now = now

		chunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "5h") {
			t.Errorf("rendered text = %q, want it to contain the 5h countdown", chunkText(chunks))
		}
		last := chunks[len(chunks)-1]
		if last.FG != rc.Theme.Muted {
			t.Errorf("reset chunk FG = %+v, want theme.Muted %+v", last.FG, rc.Theme.Muted)
		}
	})

	t.Run("a reset time already in the past renders now", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		resetsAt := now.Add(-1 * time.Hour)
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 10, ResetsAt: resetsAt.Unix()},
			},
		}, nil)
		rc.Now = now

		chunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "now") {
			t.Errorf("rendered text = %q, want it to contain now", chunkText(chunks))
		}
	})

	t.Run("bar cells use the fixed position gradient, matching context_window", func(t *testing.T) {
		t.Parallel()
		th := testTheme(t)
		rc := newTestContext(t, &input.Payload{
			RateLimits: &input.RateLimits{FiveHour: &input.RateLimitWindow{UsedPercentage: 100}},
		}, nil)

		chunks, ok := newRateLimitSegment(windowFiveHour).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		// chunks[0] = icon+label, chunks[1] = separator space, chunks[2:2+width] = bar cells.
		const barStart = 2
		for i := range rateLimitGaugeWidth {
			want := positionGradientColor(&th, i, rateLimitGaugeWidth)
			got := chunks[barStart+i].FG
			if got != want {
				t.Errorf("bar cell %d FG = %+v, want %+v (position gradient)", i, got, want)
			}
		}
	})
}

func TestFormatResetIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{name: "zero is now", d: 0, want: "now"},
		{name: "negative is now", d: -time.Hour, want: "now"},
		{name: "minutes only", d: 42 * time.Minute, want: "42m"},
		{name: "exact hour", d: 5 * time.Hour, want: "5h"},
		{name: "hours, minutes dropped", d: 4*time.Hour + 45*time.Minute, want: "4h"},
		{name: "exact day", d: 3 * 24 * time.Hour, want: "3d"},
		{name: "days and hours", d: 3*24*time.Hour + 2*time.Hour, want: "3d 2h"},
		{name: "days, hours dropped when zero", d: 7 * 24 * time.Hour, want: "7d"},
		{name: "rounds to the nearest minute", d: 90*time.Second + 29*time.Second, want: "2m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := formatResetIn(tt.d); got != tt.want {
				t.Errorf("formatResetIn(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
