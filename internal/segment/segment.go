package segment

import (
	"time"

	"github.com/scrothers/statusline/internal/config"
	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// RenderContext carries everything a Segment needs to render itself: the
// parsed payload, resolved config and theme, terminal width, and (for the
// git segment) pre-collected git status.
type RenderContext struct {
	Payload *input.Payload
	Config  *config.Config
	Theme   *theme.Theme
	Columns int
	Now     time.Time
	Git     *gitstatus.Status
}

// Segment renders one piece of the statusline. Render returns ok=false when
// the segment has nothing to show for the current payload (e.g. no open PR),
// in which case it's omitted entirely rather than rendered empty.
type Segment interface {
	ID() string
	// Priority ranks how long a segment survives width-pressure truncation;
	// higher values are dropped last.
	Priority() int
	Render(rc *RenderContext) (chunks []style.Chunk, ok bool)
}

// Registry returns every built-in segment keyed by ID. It's a constructor,
// not a package-level var, so there's no shared mutable global to guard.
func Registry() map[string]Segment {
	segments := []Segment{
		modelSegment{},
		providerSegment{},
		billingSegment{},
		directorySegment{},
		gitSegment{},
		contextWindowSegment{},
		costSegment{},
		durationSegment{},
		newRateLimitSegment(windowFiveHour),
		newRateLimitSegment(windowSevenDay),
		prSegment{},
		vimSegment{},
		agentSegment{},
		effortSegment{},
		outputStyleSegment{},
		thinkingSegment{},
		cacheSegment{},
		sessionNameSegment{},
		linesChangedSegment{},
		tokenCountsSegment{},
		repoSegment{},
		worktreeSegment{},
	}
	reg := make(map[string]Segment, len(segments))
	for _, s := range segments {
		reg[s.ID()] = s
	}
	return reg
}
