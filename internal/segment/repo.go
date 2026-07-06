package segment

import (
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// repoSegment renders the git remote's host, owner, and name as three
// distinctly colored pieces joined like a path (e.g.
// "github.com/scrothers/statusline"), rather than one flat string.
// Omitted outside a git repository or when there's no origin remote.
type repoSegment struct{}

func (repoSegment) ID() string { return "repo" }

func (repoSegment) Priority() int { return 85 }

func (repoSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	if rc.Payload.Workspace == nil || rc.Payload.Workspace.Repo == nil {
		return nil, false
	}
	r := rc.Payload.Workspace.Repo

	icon := theme.Glyph(theme.IconRepo, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
		{Text: " " + r.Host, FG: rc.Theme.Muted},
		{Text: "/" + r.Owner, FG: rc.Theme.TextSecondary},
		{Text: "/" + r.Name, FG: rc.Theme.IdentityText},
	}, true
}
