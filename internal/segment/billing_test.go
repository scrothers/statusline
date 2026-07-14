package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestBillingSegment(t *testing.T) {
	t.Parallel()

	t.Run("renders when cost is present and rate limits are absent", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			Cost: &input.Cost{TotalCostUSD: 1.23},
		}, nil)
		disableNerdFont(rc)

		chunks, ok := billingSegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "API") {
			t.Errorf("rendered text = %q, want it to contain the API fallback", chunkText(chunks))
		}
	})

	t.Run("omitted when rate limits are present", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			Cost: &input.Cost{TotalCostUSD: 1.23},
			RateLimits: &input.RateLimits{
				FiveHour: &input.RateLimitWindow{UsedPercentage: 10},
			},
		}, nil)
		_, ok := billingSegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false when RateLimits is present")
		}
	})

	t.Run("omitted before the first API response of the session", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		_, ok := billingSegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false when Cost is nil")
		}
	})
}
