package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
)

func TestEffortSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent effort is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (effortSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Effort")
		}
	})

	t.Run("known levels render with a distinct icon and color each", func(t *testing.T) {
		t.Parallel()
		levels := []string{"low", "medium", "high", "xhigh", "max", "ultra"}
		seenColors := make(map[style.Color]bool, len(levels))

		for _, level := range levels {
			rc := newTestContext(t, &input.Payload{Effort: &input.Effort{Level: level}}, nil)
			chunks, ok := (effortSegment{}).Render(rc)
			if !ok {
				t.Fatalf("Render() ok = false for level %q, want true", level)
			}
			if !strings.Contains(chunkText(chunks), level) {
				t.Errorf("level %q: rendered text = %q, want it to contain the level word", level, chunkText(chunks))
			}
			if chunkText(chunks) == level {
				t.Errorf("level %q: rendered text = %q, want an icon prefix too", level, chunkText(chunks))
			}

			if seenColors[chunks[0].FG] {
				t.Errorf("level %q reused a color already used by an earlier level", level)
			}
			seenColors[chunks[0].FG] = true
		}
	})

	t.Run("progresses from green toward red across low..max", func(t *testing.T) {
		t.Parallel()
		low := effortLevels["low"].color
		medium := effortLevels["medium"].color
		high := effortLevels["high"].color
		xhigh := effortLevels["xhigh"].color
		maxLevel := effortLevels["max"].color

		// Green channel should trend down and red channel up as intensity
		// rises from low to max — not an exact monotonic guarantee for
		// every channel (yellow/orange bumps green back up briefly), but
		// red should strictly climb and end near-saturated by max.
		if low.R >= xhigh.R || xhigh.R > maxLevel.R {
			t.Errorf("red channel should rise toward max: low=%d medium=%d high=%d xhigh=%d max=%d",
				low.R, medium.R, high.R, xhigh.R, maxLevel.R)
		}
		if maxLevel.R < 200 {
			t.Errorf("max color should be strongly red (R>=200), got R=%d", maxLevel.R)
		}
		if low.G < 150 {
			t.Errorf("low color should be strongly green (G>=150), got G=%d", low.G)
		}
	})

	t.Run("ultra is purple: high red and blue, low green", func(t *testing.T) {
		t.Parallel()
		ultra := effortLevels["ultra"].color
		if ultra.G >= ultra.R || ultra.G >= ultra.B {
			t.Errorf("ultra should read as purple (low green relative to red/blue), got %+v", ultra)
		}
	})

	t.Run("unknown level still renders, in identity accent", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Effort: &input.Effort{Level: "not-a-real-level"}}, nil)
		chunks, ok := (effortSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if chunkText(chunks) != "not-a-real-level" {
			t.Errorf("rendered text = %q, want the raw level string", chunkText(chunks))
		}
		if chunks[0].FG != rc.Theme.IdentityAccent {
			t.Errorf("FG = %+v, want identity accent %+v", chunks[0].FG, rc.Theme.IdentityAccent)
		}
	})
}
