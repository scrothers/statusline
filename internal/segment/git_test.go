package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/gitstatus"
	"github.com/scrothers/statusline/internal/input"
)

func TestGitSegment(t *testing.T) {
	t.Parallel()

	t.Run("nil git status is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (gitSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Git")
		}
	})

	t.Run("not a repo is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{NotARepo: true})
		if _, ok := (gitSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for NotARepo")
		}
	})

	t.Run("clean repo shows branch only", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Branch: "main"})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "main") {
			t.Errorf("rendered text = %q, want it to contain main", text)
		}
	})

	t.Run("dirty repo shows counts", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{
			Branch: "main", Staged: 2, Modified: 1, Untracked: 3, Ahead: 1, Behind: 2,
		})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		for _, want := range []string{"2", "1", "3"} {
			if !strings.Contains(text, want) {
				t.Errorf("rendered text = %q, want it to contain %q", text, want)
			}
		}
	})

	t.Run("branch name uses identity text color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Branch: "main"})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) < 2 {
			t.Fatal("Render() produced fewer than 2 chunks")
		}
		if chunks[1].FG != rc.Theme.IdentityText {
			t.Errorf("branch name FG = %+v, want theme.IdentityText %+v", chunks[1].FG, rc.Theme.IdentityText)
		}
	})

	t.Run("badge icon carries category color, count uses secondary text color", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Branch: "main", Staged: 2})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if len(chunks) != 4 {
			t.Fatalf("Render() produced %d chunks, want 4 (branch icon, branch name, staged icon, staged count)", len(chunks))
		}
		if chunks[2].FG != rc.Theme.Success {
			t.Errorf("staged icon FG = %+v, want theme.Success %+v", chunks[2].FG, rc.Theme.Success)
		}
		if chunks[3].FG != rc.Theme.TextSecondary {
			t.Errorf("staged count FG = %+v, want theme.TextSecondary %+v", chunks[3].FG, rc.Theme.TextSecondary)
		}
	})

	t.Run("conflicts render with the alert icon, not a bare ASCII bang", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Branch: "main", Conflicts: 1})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if strings.Contains(text, "!") {
			t.Errorf("rendered text = %q, want no bare ASCII ! (an icon should carry this meaning)", text)
		}
		if len(chunks) != 4 {
			t.Fatalf("Render() produced %d chunks, want 4 (branch icon, branch name, conflict icon, conflict count)", len(chunks))
		}
		if chunks[2].FG != rc.Theme.Danger {
			t.Errorf("conflict icon FG = %+v, want theme.Danger %+v", chunks[2].FG, rc.Theme.Danger)
		}
	})

	t.Run("detached head shows short oid", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Detached: true, OID: "abcdef1234567890"})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		if !strings.Contains(text, "abcdef1") {
			t.Errorf("rendered text = %q, want it to contain short oid abcdef1", text)
		}
		if strings.Contains(text, "abcdef1234567890") {
			t.Errorf("rendered text = %q, oid should be truncated to 7 chars", text)
		}
	})

	t.Run("detached head with no oid falls back to HEAD", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, &gitstatus.Status{Detached: true})
		chunks, ok := (gitSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "HEAD") {
			t.Errorf("rendered text = %q, want it to contain HEAD", chunkText(chunks))
		}
	})
}
