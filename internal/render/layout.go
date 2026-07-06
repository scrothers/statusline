package render

import (
	"sort"

	"github.com/scrothers/statusline/internal/style"
)

// neverDropPriority is the priority threshold at or above which a segment
// is never automatically dropped by fitToWidth (model and directory render
// at this priority and are only ever shrunk internally, never omitted).
const neverDropPriority = 100

// lineSegment is one rendered segment ready for width-fitting and joining:
// its chunks plus the metadata layout needs without re-deriving it.
type lineSegment struct {
	id       string
	priority int
	chunks   []style.Chunk
	bg       style.Color // uniform background across the segment's chunks

	// breakBefore forces a hard gap (taper to the terminal's default
	// background, a plain space, taper back in) before this segment even
	// when it shares a background color with the previous one — the config
	// "gap" marker between unrelated clusters on the same line (e.g. the
	// context-window gauge vs. cost+duration) sets this rather than relying
	// on a background-color difference, so breathing room never means
	// painting more background color.
	breakBefore bool
}

func chunksWidth(chunks []style.Chunk) int {
	total := 0
	for _, c := range chunks {
		total += c.Width()
	}
	return total
}

// lineWidth estimates the rendered width of a line including join glyphs:
// one cell for a pill-to-pill connector or an end cap, two for a pill/badge
// transition (a closing or opening cap plus a plain space), three for a
// " · " divider between two badges, and a breakBefore gap costs one cap per
// side that's a pill plus the space between (so 1-3 cells, symmetric with
// the transition it's forcing).
func lineWidth(segments []lineSegment) int {
	if len(segments) == 0 {
		return 0
	}
	total := 0
	for i, s := range segments {
		total += chunksWidth(s.chunks)
		if i == 0 {
			continue
		}
		prev := segments[i-1]
		switch {
		case s.breakBefore:
			if prev.bg.Valid {
				total++
			}
			total++ // the gap itself
			if s.bg.Valid {
				total++
			}
		case prev.bg.Valid && s.bg.Valid:
			total++
		case prev.bg.Valid != s.bg.Valid:
			total += 2
		default:
			total += 3
		}
	}
	if segments[0].bg.Valid {
		total++
	}
	if segments[len(segments)-1].bg.Valid {
		total++
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
