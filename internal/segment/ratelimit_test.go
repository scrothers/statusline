package segment

import (
	"strings"
	"testing"

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
		// chunks[0] = icon+label, chunks[1] = " ⟨", chunks[2:2+width] = bar cells.
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
