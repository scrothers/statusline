package gitstatus

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheKey(t *testing.T) {
	t.Parallel()

	a := cacheKey("session-1", "/home/user/repo")
	b := cacheKey("session-1", "/home/user/repo")
	if a != b {
		t.Errorf("cacheKey() is not deterministic: %q != %q", a, b)
	}
	if len(a) != 32 {
		t.Errorf("cacheKey() length = %d, want 32 (16 bytes hex-encoded)", len(a))
	}

	t.Run("differs by session", func(t *testing.T) {
		t.Parallel()
		if cacheKey("session-1", "/repo") == cacheKey("session-2", "/repo") {
			t.Error("cacheKey() collided across different session IDs")
		}
	})

	t.Run("differs by dir", func(t *testing.T) {
		t.Parallel()
		if cacheKey("session-1", "/repo-a") == cacheKey("session-1", "/repo-b") {
			t.Error("cacheKey() collided across different dirs")
		}
	})
}

func TestCachedStatus(t *testing.T) {
	t.Parallel()

	t.Run("not a repo short-circuits to NotARepo, ignoring any stale status", func(t *testing.T) {
		t.Parallel()
		entry := cachedEntry{NotARepo: true, Status: Status{Branch: "stale-branch"}}
		got := cachedStatus(entry)
		if !got.NotARepo || got.Branch != "" {
			t.Errorf("cachedStatus() = %+v, want {NotARepo: true} only", got)
		}
	})

	t.Run("returns the cached status as-is", func(t *testing.T) {
		t.Parallel()
		want := Status{Branch: "main", Staged: 2}
		got := cachedStatus(cachedEntry{Status: want})
		if got != want {
			t.Errorf("cachedStatus() = %+v, want %+v", got, want)
		}
	})
}

func TestGitCacheDir(t *testing.T) {
	t.Parallel()

	dir, err := gitCacheDir()
	if err != nil {
		t.Fatalf("gitCacheDir() error = %v", err)
	}
	if filepath.Base(dir) != "git" || filepath.Base(filepath.Dir(dir)) != "statusline" {
		t.Errorf("gitCacheDir() = %q, want a path ending in statusline/git", dir)
	}
}

func TestHeadMTimeUnixNano(t *testing.T) {
	t.Parallel()

	t.Run("no .git/HEAD returns zero", func(t *testing.T) {
		t.Parallel()
		if got := headMTimeUnixNano(t.TempDir()); got != 0 {
			t.Errorf("headMTimeUnixNano() = %d, want 0 for a directory with no .git/HEAD", got)
		}
	})

	t.Run("returns the file's actual mtime", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, ".git"), 0o750); err != nil {
			t.Fatalf("Mkdir() error = %v", err)
		}
		headPath := filepath.Join(dir, ".git", "HEAD")
		if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		info, err := os.Stat(headPath)
		if err != nil {
			t.Fatalf("Stat() error = %v", err)
		}

		got := headMTimeUnixNano(dir)
		if got != info.ModTime().UnixNano() {
			t.Errorf("headMTimeUnixNano() = %d, want %d (the file's actual mtime)", got, info.ModTime().UnixNano())
		}
	})
}
