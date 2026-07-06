package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// worktreeSegment renders the active worktree's name (falling back to its
// branch when the name is empty). Present only during --worktree sessions.
type worktreeSegment struct{}

func (worktreeSegment) ID() string { return "worktree" }

func (worktreeSegment) Priority() int { return 25 }

func (worktreeSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	wt := rc.Payload.Worktree
	if wt == nil {
		return nil, false
	}

	name := wt.Name
	if name == "" {
		name = wt.Branch
	}
	if name == "" {
		return nil, false
	}

	icon := theme.Glyph(theme.IconWorktree, rc.Config.NerdFontEnabled())
	return []style.Chunk{{Text: icon + " " + name, FG: rc.Theme.Info}}, true
}
