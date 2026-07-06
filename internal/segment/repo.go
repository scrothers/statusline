package segment

import (
	"strings"

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

	icon := theme.Glyph(repoHostIconKey(r.Host), rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
		{Text: " " + r.Host, FG: rc.Theme.Muted},
		{Text: "/" + r.Owner, FG: rc.Theme.TextSecondary},
		{Text: "/" + r.Name, FG: rc.Theme.IdentityText},
	}, true
}

// repoHostIconKey picks a host-branded icon by matching known git forges as
// a substring of the host (covers both public instances like "github.com"
// and enterprise/self-hosted ones like "github.company.com"), falling back
// to a generic git icon for anything else (e.g. a bare "git.company.com").
func repoHostIconKey(host string) string {
	h := strings.ToLower(host)
	switch {
	case strings.Contains(h, "gitlab"):
		return theme.IconRepoGitLab
	case strings.Contains(h, "github"):
		return theme.IconRepoGitHub
	case strings.Contains(h, "forgejo"):
		return theme.IconRepoForgejo
	default:
		return theme.IconRepoGit
	}
}
