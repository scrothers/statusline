package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestSessionNameSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent session name is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (sessionNameSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for empty SessionName")
		}
	})

	t.Run("renders session name", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{SessionName: "big-refactor"}, nil)
		chunks, ok := (sessionNameSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "big-refactor") {
			t.Errorf("rendered text = %q, want it to contain big-refactor", chunkText(chunks))
		}
	})
}

func TestLinesChangedSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent cost is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (linesChangedSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Cost")
		}
	})

	t.Run("zero added and removed is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{}}, nil)
		if _, ok := (linesChangedSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when both are zero")
		}
	})

	t.Run("renders both added and removed", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 342, TotalLinesRemoved: 58}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "+342") || !strings.Contains(text, "-58") {
			t.Errorf("rendered text = %q, want it to contain +342 and -58", text)
		}
	})

	t.Run("renders added only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 10}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if strings.Contains(chunkText(chunks), "-") {
			t.Errorf("rendered text = %q, should not contain a removed count", chunkText(chunks))
		}
	})

	t.Run("renders removed only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesRemoved: 5}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if strings.Contains(chunkText(chunks), "+") {
			t.Errorf("rendered text = %q, should not contain an added count", chunkText(chunks))
		}
	})
}
