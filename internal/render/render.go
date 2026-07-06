package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/scrothers/statusline/internal/segment"
	"github.com/scrothers/statusline/internal/style"
)

// Render composes every enabled line's segments (as looked up in registry)
// into the final statusline text. It never panics: a segment that panics is
// recovered per-segment (logged to stderr, then omitted) so one bad segment
// only drops itself, never the whole line or the whole render.
func Render(rc *segment.RenderContext, registry map[string]segment.Segment) string {
	lines := make([]string, 0, len(rc.Config.Lines))
	for _, lc := range rc.Config.Lines {
		if !lc.Enabled {
			continue
		}
		line := renderLine(rc, registry, lc.Segments)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GapMarker is a reserved entry in a LineConfig's Segments list: it inserts
// breathing room (a plain-space break tapering to the terminal's default
// background on both sides) before the next segment, instead of the normal
// connector — for unrelated clusters that would otherwise glue together
// silently because they happen to share a background color. It never adds
// background color; it's a real gap.
const GapMarker = "gap"

func renderLine(rc *segment.RenderContext, registry map[string]segment.Segment, ids []string) string {
	segments := make([]lineSegment, 0, len(ids))
	pendingGap := false
	for _, id := range ids {
		if id == GapMarker {
			pendingGap = true
			continue
		}
		seg, ok := registry[id]
		if !ok {
			continue
		}
		if cfg, ok := rc.Config.Segments[id]; ok && cfg.Enabled != nil && !*cfg.Enabled {
			continue
		}
		chunks := safeRender(seg, rc)
		if len(chunks) == 0 {
			continue
		}
		segments = append(segments, lineSegment{
			id: id, priority: seg.Priority(), chunks: chunks, bg: chunks[0].BG,
			breakBefore: pendingGap,
		})
		pendingGap = false
	}

	segments = fitToWidth(segments, rc.Columns)
	if len(segments) == 0 {
		return ""
	}
	return joinLine(segments, rc.Config.Separator.Style, rc.Theme.Muted)
}

// safeRender recovers a panic from an individual segment so it can't take
// down the whole line; the panic is logged to stderr (stdout is reserved
// for the rendered statusline) and the segment is treated as having nothing
// to show.
func safeRender(seg segment.Segment, rc *segment.RenderContext) (chunks []style.Chunk) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "statusline: segment %q panicked: %v\n", seg.ID(), r)
			chunks = nil
		}
	}()
	c, ok := seg.Render(rc)
	if !ok {
		return nil
	}
	return c
}
