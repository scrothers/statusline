package segment

import (
	"os"
	"strings"
	"unicode/utf8"

	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// directoryMaxLen bounds the breadcrumb-truncated directory display.
const directoryMaxLen = 32

// directorySegment renders the current working directory as a
// breadcrumb-truncated path.
type directorySegment struct{}

func (directorySegment) ID() string { return "directory" }

// Priority matches modelSegment: directory is never fully dropped, only
// shrunk via its own breadcrumb truncation.
func (directorySegment) Priority() int { return 100 }

func (directorySegment) Render(rc *RenderContext) ([]style.Chunk, bool) {
	dir := rc.Payload.CWD
	if rc.Payload.Workspace != nil && rc.Payload.Workspace.CurrentDir != "" {
		dir = rc.Payload.Workspace.CurrentDir
	}
	if dir == "" {
		return nil, false
	}

	home, _ := os.UserHomeDir()
	text := breadcrumb(dir, home, directoryMaxLen)

	icon := theme.Glyph(theme.IconDirectory, rc.Config.NerdFontEnabled())
	return []style.Chunk{
		{Text: icon, FG: rc.Theme.IdentityAccent},
		{Text: " " + text, FG: rc.Theme.IdentityText},
	}, true
}

// breadcrumb shortens path for display: it substitutes home with "~", then
// (only if still too long) shrinks every middle segment to its first rune,
// then falls back to the last two segments with a leading ellipsis, then
// finally the last segment alone.
func breadcrumb(path, home string, maxLen int) string {
	if home != "" && (path == home || strings.HasPrefix(path, home+"/")) {
		path = "~" + strings.TrimPrefix(path, home)
	}
	if utf8.RuneCountInString(path) <= maxLen {
		return path
	}

	parts := strings.Split(path, "/")
	shrunk := make([]string, len(parts))
	copy(shrunk, parts)
	for i := 1; i < len(shrunk)-1; i++ {
		if shrunk[i] == "" {
			continue
		}
		r := []rune(shrunk[i])
		shrunk[i] = string(r[0])
	}
	if candidate := strings.Join(shrunk, "/"); utf8.RuneCountInString(candidate) <= maxLen {
		return candidate
	}

	if len(parts) >= 2 {
		candidate := "…/" + strings.Join(parts[len(parts)-2:], "/")
		if utf8.RuneCountInString(candidate) <= maxLen {
			return candidate
		}
	}
	return parts[len(parts)-1]
}
