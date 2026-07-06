package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
)

func TestRepoSegment(t *testing.T) {
	t.Parallel()

	t.Run("absent workspace is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		if _, ok := (repoSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Workspace")
		}
	})

	t.Run("absent repo is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{Workspace: &input.Workspace{}}, nil)
		if _, ok := (repoSegment{}).Render(rc); ok {
			t.Error("Render() ok = true, want false for nil Repo")
		}
	})

	t.Run("renders host, owner, and name distinctly", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			Workspace: &input.Workspace{Repo: &input.Repo{Host: "github.com", Owner: "scrothers", Name: "statusline"}},
		}, nil)
		chunks, ok := (repoSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		text := chunkText(chunks)
		for _, want := range []string{"github.com", "scrothers", "statusline"} {
			if !strings.Contains(text, want) {
				t.Errorf("rendered text = %q, want it to contain %q", text, want)
			}
		}
		// Each field renders as its own chunk (a real "breakdown"), not one
		// flat string.
		if len(chunks) < 3 {
			t.Errorf("chunks = %+v, want at least 3 (host, owner, name broken out)", chunks)
		}
	})
}
