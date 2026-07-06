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

	t.Run("renders both added and removed with no ASCII sign", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 342, TotalLinesRemoved: 58}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "342") || !strings.Contains(text, "58") {
			t.Errorf("rendered text = %q, want it to contain 342 and 58", text)
		}
		if strings.ContainsAny(text, "+-") {
			t.Errorf("rendered text = %q, want no ASCII +/- (the icons alone carry that meaning)", text)
		}
	})

	t.Run("leads with the pencil icon in warning color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 342, TotalLinesRemoved: 58}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
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

	t.Run("pencil icon has two trailing spaces for breathing room", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 342, TotalLinesRemoved: 58}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) == 0 || !strings.HasSuffix(chunks[0].Text, "  ") {
			t.Errorf("pencil chunk = %q, want it to end with two spaces", chunks[0].Text)
		}
	})

	t.Run("diff icons carry semantic color, counts use secondary text color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 342, TotalLinesRemoved: 58}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) != 5 {
			t.Fatalf("Render() produced %d chunks, want 5 (pencil, +icon, +count, -icon, -count)", len(chunks))
		}
		if chunks[1].FG != rc.Theme.Success {
			t.Errorf("added icon FG = %+v, want theme.Success %+v", chunks[1].FG, rc.Theme.Success)
		}
		if chunks[2].FG != rc.Theme.TextSecondary {
			t.Errorf("added count FG = %+v, want theme.TextSecondary %+v", chunks[2].FG, rc.Theme.TextSecondary)
		}
		if chunks[3].FG != rc.Theme.Danger {
			t.Errorf("removed icon FG = %+v, want theme.Danger %+v", chunks[3].FG, rc.Theme.Danger)
		}
		if chunks[4].FG != rc.Theme.TextSecondary {
			t.Errorf("removed count FG = %+v, want theme.TextSecondary %+v", chunks[4].FG, rc.Theme.TextSecondary)
		}
	})

	t.Run("renders added only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesAdded: 10}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "10") {
			t.Errorf("rendered text = %q, want it to contain 10", chunkText(chunks))
		}
	})

	t.Run("renders removed only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Cost: &input.Cost{TotalLinesRemoved: 5}}, nil)
		chunks, ok := (linesChangedSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "5") {
			t.Errorf("rendered text = %q, want it to contain 5", chunkText(chunks))
		}
	})
}
