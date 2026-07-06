package segment

import (
	"strings"
	"testing"

	"github.com/scrothers/statusline/internal/input"
	"github.com/scrothers/statusline/internal/theme"
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

	t.Run("repo name uses identity text color, matching other headline labels", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			Workspace: &input.Workspace{Repo: &input.Repo{Host: "github.com", Owner: "scrothers", Name: "statusline"}},
		}, nil)
		chunks, ok := (repoSegment{}).Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		last := chunks[len(chunks)-1]
		if last.FG != rc.Theme.IdentityText {
			t.Errorf("repo name FG = %+v, want theme.IdentityText %+v", last.FG, rc.Theme.IdentityText)
		}
	})
}

func TestRepoHostIconKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		host string
		want string
	}{
		{host: "github.com", want: theme.IconRepoGitHub},
		{host: "github.enterprise.example.com", want: theme.IconRepoGitHub},
		{host: "gitlab.com", want: theme.IconRepoGitLab},
		{host: "gitlab.example.com", want: theme.IconRepoGitLab},
		{host: "codeberg.org", want: theme.IconRepoGit}, // forgejo-based but host doesn't say so
		{host: "forgejo.example.com", want: theme.IconRepoForgejo},
		{host: "git.example.com", want: theme.IconRepoGit},
		{host: "", want: theme.IconRepoGit},
	}
	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			t.Parallel()
			if got := repoHostIconKey(tt.host); got != tt.want {
				t.Errorf("repoHostIconKey(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}
