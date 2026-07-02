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
}
