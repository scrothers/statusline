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
}
