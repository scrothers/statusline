package render

import (
	"sort"

	"github.com/scrothers/statusline/internal/style"
)

// neverDropPriority is the priority threshold at or above which a segment
// is never automatically dropped by fitToWidth (model and directory render
// at this priority and are only ever shrunk internally, never omitted).
const neverDropPriority = 100

// dividerWidth is the display width of dividerText (defined in divider.go)
// between two adjacent segments on a line: 3 spaces + a single-width glyph
// + 3 spaces.
const dividerWidth = 7

// lineSegment is one rendered segment ready for width-fitting and joining:
// its chunks plus the metadata layout needs without re-deriving it.
type lineSegment struct {
	id       string
	priority int
	chunks   []style.Chunk
}

func chunksWidth(chunks []style.Chunk) int {
	total := 0
	for _, c := range chunks {
		total += c.Width()
	}
	return total
}

// lineWidth estimates the rendered width of a line: every segment's content
// plus one divider between each adjacent pair.
func lineWidth(segments []lineSegment) int {
	if len(segments) == 0 {
		return 0
	}
	total := 0
	for i, s := range segments {
		total += chunksWidth(s.chunks)
		if i > 0 {
			total += dividerWidth
		}
	}
	return total
}

// fitToWidth drops segments in ascending priority order (lowest first) until
// the line fits columns, but never drops a segment at or above
// neverDropPriority — model and directory always render, however tight the
// terminal, since they self-truncate instead.
func fitToWidth(segments []lineSegment, columns int) []lineSegment {
	if columns <= 0 || lineWidth(segments) <= columns {
		return segments
	}

	order := make([]int, len(segments))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		return segments[order[a]].priority < segments[order[b]].priority
	})

	dropped := make(map[int]bool, len(segments))
	for _, idx := range order {
		if segments[idx].priority >= neverDropPriority {
			break // ascending order: everything after this is also protected
		}
		dropped[idx] = true
		remaining := keep(segments, dropped)
		if lineWidth(remaining) <= columns {
			return remaining
		}
	}
	return keep(segments, dropped)
}

func keep(segments []lineSegment, dropped map[int]bool) []lineSegment {
	out := make([]lineSegment, 0, len(segments))
	for i, s := range segments {
		if !dropped[i] {
			out = append(out, s)
		}
	}
	return out
}
