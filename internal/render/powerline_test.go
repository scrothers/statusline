package render

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/style"
)

func TestCapsFor(t *testing.T) {
	t.Parallel()

	if got := capsFor("hard"); got != capStyles["hard"] {
		t.Errorf("capsFor(hard) = %+v, want %+v", got, capStyles["hard"])
	}
	if got := capsFor("not-a-style"); got != capStyles[defaultCapStyle] {
		t.Errorf("capsFor(unknown) = %+v, want default %+v", got, capStyles[defaultCapStyle])
	}
	if got := capsFor(""); got != capStyles[defaultCapStyle] {
		t.Errorf("capsFor(empty) = %+v, want default %+v", got, capStyles[defaultCapStyle])
	}
}

// TestJoinLine mutates the NO_COLOR env var to isolate text content from
// ANSI codes, so neither it nor its subtests run in parallel.
func TestJoinLine(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	muted := style.RGB(1, 1, 1)

	t.Run("empty input", func(t *testing.T) {
		if got := joinLine(nil, "rounded", muted); got != "" {
			t.Errorf("joinLine(nil) = %q, want empty", got)
		}
	})

	t.Run("single pill contains its text", func(t *testing.T) {
		segs := []lineSegment{pillSegment("model", 100, " Opus ", style.RGB(1, 2, 3))}
		got := joinLine(segs, "rounded", muted)
		if !strings.Contains(got, "Opus") {
			t.Errorf("joinLine() = %q, want it to contain Opus", got)
		}
	})

	t.Run("badge to badge uses a dot divider", func(t *testing.T) {
		segs := []lineSegment{badgeSegment("vim", 30, "INSERT"), badgeSegment("agent", 30, "reviewer")}
		got := joinLine(segs, "rounded", muted)
		if !strings.Contains(got, "·") {
			t.Errorf("joinLine() = %q, want a · divider between badges", got)
		}
		if !strings.Contains(got, "INSERT") || !strings.Contains(got, "reviewer") {
			t.Errorf("joinLine() = %q, missing badge content", got)
		}
	})

	t.Run("all segment text survives regardless of style", func(t *testing.T) {
		segs := []lineSegment{
			pillSegment("model", 100, "Opus", style.RGB(1, 2, 3)),
			pillSegment("directory", 100, "proj", style.RGB(4, 5, 6)),
			badgeSegment("vim", 30, "NORMAL"),
		}
		for _, capStyle := range []string{"rounded", "hard", "unknown"} {
			got := joinLine(segs, capStyle, muted)
			for _, want := range []string{"Opus", "proj", "NORMAL"} {
				if !strings.Contains(got, want) {
					t.Errorf("joinLine(%q) = %q, missing %q", capStyle, got, want)
				}
			}
		}
	})
}
