package segment

import (
	"testing"

	"github.com/scrothers/statusline/internal/config"
	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/style"
	"github.com/scrothers/statusline/internal/theme"
)

// testTheme resolves the default (Gruvbox) built-in theme for use across
// every segment test in this package.
func testTheme(t *testing.T) theme.Theme {
	t.Helper()
	registry, err := theme.LoadRegistry()
	if err != nil {
		t.Fatalf("theme.LoadRegistry() error = %v", err)
	}
	th, _ := theme.Resolve(registry, theme.DefaultName)
	return th
}

// newTestContext builds a RenderContext with default config/theme for a
// given payload and (optional) git status.
func newTestContext(t *testing.T, payload *input.Payload, git *gitstatus.Status) *RenderContext {
	t.Helper()
	th := testTheme(t)
	cfg := config.Default()
	return &RenderContext{
		Payload: payload,
		Config:  &cfg,
		Theme:   &th,
		Columns: 120,
		Git:     git,
	}
}

// chunkText concatenates every chunk's text, for tests that only care about
// the rendered content, not per-chunk coloring.
func chunkText(chunks []style.Chunk) string {
	var s string
	for _, c := range chunks {
		s += c.Text
	}
	return s
}

func TestRegistry(t *testing.T) {
	t.Parallel()

	reg := Registry()
	wantIDs := []string{
		"model", "provider", "directory", "git", "context_window", "cost", "duration",
		"ratelimit_5h", "ratelimit_7d", "pr", "vim", "agent", "effort", "output_style",
		"thinking", "cache", "session_name", "lines_changed", "token_counts", "repo", "worktree",
	}
	if len(reg) != len(wantIDs) {
		t.Fatalf("Registry() has %d segments, want %d", len(reg), len(wantIDs))
	}
	for _, id := range wantIDs {
		seg, ok := reg[id]
		if !ok {
			t.Errorf("Registry() missing segment %q", id)
			continue
		}
		if seg.ID() != id {
			t.Errorf("Registry()[%q].ID() = %q, want %q", id, seg.ID(), id)
		}
	}
}

// TestSegmentPriorities pins every segment's declared drop priority. This is
// a real regression test, not just a coverage exercise: fitToWidth's
// truncation order depends entirely on these values, so an accidental
// change here would silently reorder which segments survive a narrow
// terminal.
func TestSegmentPriorities(t *testing.T) {
	t.Parallel()

	want := map[string]int{
		"model":          100,
		"provider":       15,
		"directory":      100,
		"git":            90,
		"context_window": 90,
		"repo":           85,
		"cost":           70,
		"duration":       65,
		"pr":             60,
		"lines_changed":  55,
		"ratelimit_5h":   50,
		"ratelimit_7d":   50,
		"token_counts":   50,
		"session_name":   45,
		"effort":         40,
		"cache":          35,
		"vim":            30,
		"agent":          30,
		"thinking":       25,
		"worktree":       25,
		"output_style":   20,
	}

	reg := Registry()
	if len(want) != len(reg) {
		t.Fatalf("test table has %d entries, Registry() has %d; keep them in sync", len(want), len(reg))
	}
	for id, wantPriority := range want {
		seg, ok := reg[id]
		if !ok {
			t.Errorf("Registry() missing segment %q", id)
			continue
		}
		if got := seg.Priority(); got != wantPriority {
			t.Errorf("Registry()[%q].Priority() = %d, want %d", id, got, wantPriority)
		}
	}
}
