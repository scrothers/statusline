package render

import (
	"testing"

	"github.com/scrothers/statusline/internal/style"
)

func pillSegment(id string, priority int, text string, bg style.Color) lineSegment {
	return lineSegment{id: id, priority: priority, chunks: []style.Chunk{{Text: text, BG: bg}}, bg: bg}
}

func badgeSegment(id string, priority int, text string) lineSegment {
	return lineSegment{id: id, priority: priority, chunks: []style.Chunk{{Text: text, BG: style.Default}}, bg: style.Default}
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
	bg := style.RGB(1, 2, 3)

	tests := []struct {
		name string
		segs []lineSegment
		want int
	}{
		{name: "empty", segs: nil, want: 0},
		{
			name: "single pill: content + open + close cap",
			segs: []lineSegment{pillSegment("a", 100, "abc", bg)},
			want: 3 + 2,
		},
		{
			name: "two pills: content + 2 caps + 1 connector",
			segs: []lineSegment{pillSegment("a", 100, "ab", bg), pillSegment("b", 90, "cd", bg)},
			want: 2 + 2 + 2 + 1,
		},
		{
			name: "pill then badge: content + open cap + close cap + space, no trailing cap",
			segs: []lineSegment{pillSegment("a", 100, "ab", bg), badgeSegment("b", 30, "cd")},
			want: 2 + 2 + 1 /*open cap*/ + 2, /*close cap + space*/
		},
		{
			name: "two badges: content + divider, no caps",
			segs: []lineSegment{badgeSegment("a", 30, "ab"), badgeSegment("b", 30, "cd")},
			want: 2 + 2 + 3,
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
	bg := style.RGB(1, 2, 3)

	t.Run("fits already, unchanged", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{pillSegment("model", 100, "abc", bg)}
		got := fitToWidth(segs, 80)
		if len(got) != 1 {
			t.Fatalf("fitToWidth() = %v, want unchanged", got)
		}
	})

	t.Run("zero columns means no truncation", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{pillSegment("model", 100, "a very long piece of text indeed", bg)}
		got := fitToWidth(segs, 0)
		if len(got) != 1 {
			t.Fatalf("fitToWidth(0) dropped segments, want no-op")
		}
	})

	t.Run("drops lowest priority first", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			pillSegment("model", 100, "0123456789", bg),
			badgeSegment("vim", 30, "0123456789"),
			pillSegment("cost", 60, "0123456789", bg),
		}
		// Full width ~= 10*3 + caps ≈ way more than 15; only "model" (priority
		// 100, protected) should remain once everything droppable is dropped.
		got := fitToWidth(segs, 15)
		if len(got) != 1 || got[0].id != "model" {
			t.Fatalf("fitToWidth() = %v, want only model", got)
		}
	})

	t.Run("never drops priority >= 100", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			pillSegment("model", 100, "0123456789", bg),
			pillSegment("directory", 100, "0123456789", bg),
		}
		got := fitToWidth(segs, 1) // impossibly narrow
		if len(got) != 2 {
			t.Fatalf("fitToWidth() dropped a protected segment: %v", got)
		}
	})

	t.Run("drops just enough to fit", func(t *testing.T) {
		t.Parallel()
		segs := []lineSegment{
			pillSegment("model", 100, "model", bg),
			badgeSegment("output_style", 20, "0123456789012345"),
			badgeSegment("vim", 30, "0123456789012345"),
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
