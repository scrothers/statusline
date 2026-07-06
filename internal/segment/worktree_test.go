package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestWorktreeSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent worktree is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (worktreeSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Worktree")
		}
	})

	t.Run("renders name", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Worktree: &input.Worktree{Name: "my-feature"}}, nil)
		chunks, ok := (worktreeSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "my-feature") {
			t.Errorf("rendered text = %q, want it to contain my-feature", chunkText(chunks))
		}
	})

	t.Run("falls back to branch when name is empty", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Worktree: &input.Worktree{Branch: "worktree-my-feature"}}, nil)
		chunks, ok := (worktreeSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "worktree-my-feature") {
			t.Errorf("rendered text = %q, want it to contain the branch name", chunkText(chunks))
		}
	})

	t.Run("empty name and branch is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Worktree: &input.Worktree{}}, nil)
		if _, ok := (worktreeSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false when both Name and Branch are empty")
		}
	})

	t.Run("icon and name use the identity fact color pairing", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Worktree: &input.Worktree{Name: "my-feature"}}, nil)
		chunks, ok := (worktreeSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) != 2 {
			t.Fatalf("Render() produced %d chunks, want 2 (icon, name)", len(chunks))
		}
		if chunks[0].FG != rc.Theme.IdentityAccent {
			t.Errorf("icon FG = %+v, want theme.IdentityAccent %+v", chunks[0].FG, rc.Theme.IdentityAccent)
		}
		if chunks[1].FG != rc.Theme.IdentityText {
			t.Errorf("name FG = %+v, want theme.IdentityText %+v", chunks[1].FG, rc.Theme.IdentityText)
		}
	})
}
