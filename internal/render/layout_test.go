package render

import (
	"testing"

	"github.com/scrothers/statusline/internal/style"
)

func seg(id string, priority int, text string) lineSegment {
	return lineSegment{id: id, priority: priority, chunks: []style.Chunk{{Text: text}}}
}

func TestChunksWidth(t *testing.T) {
	t.Parallel()
	chunks := []style.Chunk{{Text: "abc"}, {Text: "def"}}
	if got := chunksWidth(chunks); got != 6 {
		t.Errorf("chunksWidth() = %d, want 6", got)
	}
}

func TestLineWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		segs []lineSegment
		want int
	}{
		{name: "empty", segs: nil, want: 0},
		{name: "single segment: just its content", segs: []lineSegment{seg("a", 100, "abc")}, want: 3},
		{
			name: "two segments: content + one divider",
			segs: []lineSegment{seg("a", 100, "ab"), seg("b", 90, "cd")},
			want: 2 + 2 + dividerWidth,
		},
		{
			name: "three segments: content + two dividers",
			segs: []lineSegment{seg("a", 100, "ab"), seg("b", 90, "cd"), seg("c", 80, "ef")},
			want: 2 + 2 + 2 + 2*dividerWidth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := lineWidth(tt.segs); got != tt.want {
				t.Errorf("lineWidth() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFitToWidth(t *testing.T) {
	t.Parallel()

	t.Run("fits already, unchanged", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{seg("model", 100, "abc")}
		got := fitToWidth(segs, 80)
		if len(got) != 1 {
			t.Fatalf("fitToWidth() = %v, want unchanged", got)
		}
	})

	t.Run("zero columns means no truncation", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{seg("model", 100, "a very long piece of text indeed")}
		got := fitToWidth(segs, 0)
		if len(got) != 1 {
			t.Fatalf("fitToWidth(0) dropped segments, want no-op")
		}
	})

	t.Run("drops lowest priority first", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			seg("model", 100, "0123456789"),
			seg("vim", 30, "0123456789"),
			seg("cost", 60, "0123456789"),
		}
		// Full width ~= 10*3 + dividers ≈ way more than 15; only "model"
		// (priority 100, protected) should remain once everything droppable
		// is dropped.
		got := fitToWidth(segs, 15)
		if len(got) != 1 || got[0].id != "model" {
			t.Fatalf("fitToWidth() = %v, want only model", got)
		}
	})

	t.Run("never drops priority >= 100", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			seg("model", 100, "0123456789"),
			seg("directory", 100, "0123456789"),
		}
		got := fitToWidth(segs, 1) // impossibly narrow
		if len(got) != 2 {
			t.Fatalf("fitToWidth() dropped a protected segment: %v", got)
		}
	})

	t.Run("drops just enough to fit", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			seg("model", 100, "model"),
			seg("output_style", 20, "0123456789012345"),
			seg("vim", 30, "0123456789012345"),
		}
		full := lineWidth(segs)
		got := fitToWidth(segs, full-1) // force exactly one drop
		if len(got) != 2 {
			t.Fatalf("fitToWidth() = %v, want exactly one segment dropped", got)
		}
		for _, s := range got {
			if s.id == "output_style" {
				t.Error("expected output_style (lowest priority) to be dropped first")
			}
		}
	})
}
