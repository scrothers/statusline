package segment

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/scrothers/statusline/internal/input"
)

func TestBreadcrumb(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		path   string
		home   string
		maxLen int
		want   string
	}{
		{name: "short path unchanged", path: "/home/user/proj", home: "/home/user", maxLen: 32, want: "~/proj"},
		{name: "home directory itself", path: "/home/user", home: "/home/user", maxLen: 32, want: "~"},
		{name: "no home match leaves path alone", path: "/var/lib/data", home: "/home/user", maxLen: 32, want: "/var/lib/data"},
		{
			// Every middle segment ("code", "some-long-project-name",
			// "internal") shrinks to its first rune.
			name:   "long path shrinks middle segments",
			path:   "/home/user/code/some-long-project-name/internal/render",
			home:   "/home/user",
			maxLen: 30,
			want:   "~/c/s/i/render",
		},
		{
			name:   "extreme case falls back to last segment only",
			path:   "/a-single-extremely-long-path-segment-that-cannot-shrink-at-all",
			home:   "",
			maxLen: 5,
			want:   "a-single-extremely-long-path-segment-that-cannot-shrink-at-all",
		},
		{
			name:   "no separators skips straight to the last-segment fallback",
			path:   "justonesegmentwithnoslashesatall",
			home:   "",
			maxLen: 5,
			want:   "justonesegmentwithnoslashesatall",
		},
		{
			name:   "empty middle segment from a double slash is skipped, not shrunk",
			path:   "/home//verylongmiddlesegmentnamehere/end",
			home:   "",
			maxLen: 10,
			want:   "/h//v/end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := breadcrumb(tt.path, tt.home, tt.maxLen); got != tt.want {
				t.Errorf("breadcrumb(%q, %q, %d) = %q, want %q", tt.path, tt.home, tt.maxLen, got, tt.want)
			}
		})
	}
}

// TestBreadcrumb_fallsBackToLastTwoSegments exercises the middle tier that
// exact-string cases above don't reach: many short middle segments make the
// shrunk-to-one-rune-each candidate longer than keeping the (short) last two
// segment names in full, since shrinking's cost scales with segment count
// while the last-two-segments candidate's cost scales with those two names'
// length instead.
func TestBreadcrumb_fallsBackToLastTwoSegments(t *testing.T) {
	t.Parallel()

	fillers := make([]string, 19, 21) // +2 capacity for the append below
	for i := range fillers {
		fillers[i] = "x"
	}
	secondLast := strings.Repeat("Y", 20)
	last := "z"
	path := "/" + strings.Join(append(fillers, secondLast, last), "/")

	got := breadcrumb(path, "", 30)
	want := "…/" + secondLast + "/" + last
	if got != want {
		t.Fatalf("breadcrumb() = %q, want %q", got, want)
	}
	if utf8.RuneCountInString(got) > 30 {
		t.Errorf("breadcrumb() length = %d, want <= 30", utf8.RuneCountInString(got))
	}
}

func TestDirectorySegment(t *testing.T) {
	t.Parallel()

	t.Run("renders current dir from workspace", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{
			Workspace: &input.Workspace{CurrentDir: "/home/user/project"},
		}, nil)
		chunks, ok := directorySegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "project") {
			t.Errorf("rendered text = %q, want it to contain project", chunkText(chunks))
		}
	})

	t.Run("falls back to top-level cwd", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{CWD: "/tmp/scratch"}, nil)
		chunks, ok := directorySegment{}.Render(rc)
		if !ok {
			t.Fatal("Render() ok = false, want true")
		}
		if !strings.Contains(chunkText(chunks), "scratch") {
			t.Errorf("rendered text = %q, want it to contain scratch", chunkText(chunks))
		}
	})

	t.Run("no directory anywhere is omitted", func(t *testing.T) {
		t.Parallel()
		rc := newTestContext(t, &input.Payload{}, nil)
		_, ok := directorySegment{}.Render(rc)
		if ok {
			t.Error("Render() ok = true, want false")
		}
	})
}
