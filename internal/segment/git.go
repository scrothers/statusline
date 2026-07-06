package segment

import (
	"strconv"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// gitSegment renders the current branch and working-tree status.
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
		{Text: " " + branchName, FG: rc.Theme.TextPrimary},
	}

	if st.Staged > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " " + theme.Glyph(theme.IconGitStaged, nerd) + " " + strconv.Itoa(st.Staged),
			FG:   rc.Theme.Success,
		})
	}
	if st.Modified > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " " + theme.Glyph(theme.IconGitModified, nerd) + " " + strconv.Itoa(st.Modified),
			FG:   rc.Theme.Warning,
		})
	}
	if st.Untracked > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " " + theme.Glyph(theme.IconGitUntracked, nerd) + " " + strconv.Itoa(st.Untracked),
			FG:   rc.Theme.Muted,
		})
	}
	if st.Conflicts > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " ! " + strconv.Itoa(st.Conflicts),
			FG:   rc.Theme.Danger,
		})
	}
	if st.Ahead > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " " + theme.Glyph(theme.IconGitAhead, nerd) + " " + strconv.Itoa(st.Ahead),
			FG:   rc.Theme.Success,
		})
	}
	if st.Behind > 0 {
		chunks = append(chunks, style.Chunk{
			Text: " " + theme.Glyph(theme.IconGitBehind, nerd) + " " + strconv.Itoa(st.Behind),
			FG:   rc.Theme.Warning,
		})
	}

	return chunks, true
}
