package segment

import (
	"strconv"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// gitSegment renders the current branch and working-tree status. Each
// status badge's icon carries its category color; the count itself renders
// in the theme's secondary text color, matching the icon/text split used by
// lines_changed and token_counts.
type gitSegment struct{}

func (gitSegment) ID() string { return "git" }

func (gitSegment) Priority() int { return 90 }

func (gitSegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	st := rc.Git
	if st == nil || st.NotARepo {
		return nil, false
	}

	nerd := rc.Config.NerdFontEnabled()

	branchColor := rc.Theme.Success
	if !st.Clean() {
		branchColor = rc.Theme.Warning
	}

	branchName := st.Branch
	if st.Detached {
		oid := st.OID
		if len(oid) > 7 {
			oid = oid[:7]
		}
		if oid == "" {
			oid = "HEAD"
		}
		branchName = oid
	}

	chunks := []style.Chunk{
		{Text: theme.Glyph(theme.IconGitBranch, nerd), FG: branchColor},
		{Text: " " + branchName, FG: rc.Theme.IdentityText},
	}

	addBadge := func(icon string, count int, color style.Color) {
		chunks = append(chunks,
			style.Chunk{Text: " " + icon, FG: color},
			style.Chunk{Text: " " + strconv.Itoa(count), FG: rc.Theme.TextSecondary},
		)
	}

	if st.Staged > 0 {
		addBadge(theme.Glyph(theme.IconGitStaged, nerd), st.Staged, rc.Theme.Success)
	}
	if st.Modified > 0 {
		// fa-pencil (also used by line 2's edit pencil) renders wide in real
		// Nerd Fonts, so it gets an extra trailing space.
		addBadge(theme.Glyph(theme.IconGitModified, nerd)+" ", st.Modified, rc.Theme.Warning)
	}
	if st.Untracked > 0 {
		addBadge(theme.Glyph(theme.IconGitUntracked, nerd), st.Untracked, rc.Theme.Muted)
	}
	if st.Conflicts > 0 {
		addBadge(theme.Glyph(theme.IconContextAlert, nerd), st.Conflicts, rc.Theme.Danger)
	}
	if st.Ahead > 0 {
		addBadge(theme.Glyph(theme.IconGitAhead, nerd), st.Ahead, rc.Theme.Success)
	}
	if st.Behind > 0 {
		addBadge(theme.Glyph(theme.IconGitBehind, nerd), st.Behind, rc.Theme.Warning)
	}

	return chunks, true
}
