//go:build integration

package gitstatus

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/scrothers/statusline/internal/config"
)

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	run("add", "file.txt")
	run("commit", "-q", "-m", "initial")
	return dir
}

func testGitConfig() config.GitConfig {
	return config.GitConfig{TimeoutMS: 2000, CacheTTLMS: 2000}
}

func TestCollect_cleanRepo(t *testing.T) {
	t.Parallel()
	dir := initRepo(t)

	st, err := Collect(context.Background(), testGitConfig(), "session-a", dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if st.Branch != "main" {
		t.Errorf("Branch = %q, want main", st.Branch)
	}
	if !st.Clean() {
		t.Errorf("Clean() = false, want true: %+v", st)
	}
}

func TestCollect_notARepo(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	st, err := Collect(context.Background(), testGitConfig(), "session-b", dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if !st.NotARepo {
		t.Errorf("NotARepo = false, want true")
	}
}

func TestCollect_dirtyRepoDetected(t *testing.T) {
	t.Parallel()
	dir := initRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	st, err := Collect(context.Background(), testGitConfig(), "session-c", dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if st.Untracked != 1 {
		t.Errorf("Untracked = %d, want 1", st.Untracked)
	}
}

func TestCollect_cacheHitAvoidsReRunningGit(t *testing.T) {
	dir := initRepo(t)
	cfg := testGitConfig()
	sessionID := "session-cache-hit"

	first, err := Collect(context.Background(), cfg, sessionID, dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if first.Untracked != 0 {
		t.Fatalf("Untracked = %d, want 0", first.Untracked)
	}

	// Add an untracked file *without* touching .git/HEAD; within the TTL a
	// cache hit should still report the pre-write (stale) state.
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	second, err := Collect(context.Background(), cfg, sessionID, dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if second.Untracked != 0 {
		t.Errorf("Untracked = %d after cache-hit window, want 0 (stale cache reused)", second.Untracked)
	}
}

func TestCollect_cacheExpiresAfterTTL(t *testing.T) {
	dir := initRepo(t)
	cfg := config.GitConfig{TimeoutMS: 2000, CacheTTLMS: 50}
	sessionID := "session-cache-expiry"

	if _, err := Collect(context.Background(), cfg, sessionID, dir); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	got, err := Collect(context.Background(), cfg, sessionID, dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if got.Untracked != 1 {
		t.Errorf("Untracked = %d after TTL expiry, want 1 (re-collected)", got.Untracked)
	}
}

func TestCollect_differentSessionsDoNotShareCache(t *testing.T) {
	t.Parallel()
	dir := initRepo(t)
	cfg := testGitConfig()

	if _, err := Collect(context.Background(), cfg, "session-x", dir); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "new.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := Collect(context.Background(), cfg, "session-y", dir)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if got.Untracked != 1 {
		t.Errorf("Untracked = %d for a different session_id, want 1 (no shared cache)", got.Untracked)
	}
}
