package gitstatus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/scrothers/statusline/internal/cache"
	"github.com/scrothers/statusline/internal/config"
)

// cachedEntry is the on-disk cache record: the parsed Status plus enough
// metadata to judge freshness without re-running git.
type cachedEntry struct {
	CapturedAtUnixNano int64  `json:"captured_at_unix_nano"`
	HeadMTimeUnixNano  int64  `json:"head_mtime_unix_nano"`
	NotARepo           bool   `json:"not_a_repo"`
	Status             Status `json:"status"`
}

// Collect returns the git status for dir. Results are cached per
// (sessionID, dir) so a burst of statusline refreshes within cfg.CacheTTLMS
// (or while .git/HEAD's mtime hasn't changed) reuses one git call instead of
// re-running git on every trigger. On a timeout or transient git error, it
// falls back to the last cached value — even if stale — rather than
// dropping the git segment outright.
func Collect(ctx context.Context, cfg config.GitConfig, sessionID, dir string) (Status, error) {
	dir = filepath.Clean(dir)
	key := cacheKey(sessionID, dir)
	cacheDir, cacheDirErr := gitCacheDir()

	var cached cachedEntry
	var hit bool
	if cacheDirErr == nil {
		cached, hit = cache.Load[cachedEntry](cacheDir, key)
	}

	headMTime := headMTimeUnixNano(dir)
	if hit {
		ttl := time.Duration(cfg.CacheTTLMS) * time.Millisecond
		fresh := time.Since(time.Unix(0, cached.CapturedAtUnixNano)) < ttl
		unchanged := headMTime != 0 && headMTime == cached.HeadMTimeUnixNano
		// Both signals must hold: TTL alone would keep serving a stale
		// branch/status through a checkout that happened seconds ago, and
		// "unchanged HEAD" alone would never expire a cache entry for a
		// repo where HEAD's mtime doesn't reflect every real change.
		if cached.NotARepo || (fresh && unchanged) {
			return cachedStatus(cached), nil
		}
	}

	status, notARepo, err := runGitStatus(ctx, cfg, dir)
	if err != nil {
		if hit {
			return cachedStatus(cached), nil
		}
		return Status{}, err
	}

	if cacheDirErr == nil {
		entry := cachedEntry{
			CapturedAtUnixNano: time.Now().UnixNano(),
			HeadMTimeUnixNano:  headMTime,
			NotARepo:           notARepo,
			Status:             status,
		}
		_ = cache.StoreAtomic(cacheDir, key, entry) // best-effort: a cache write failure shouldn't fail Collect
	}

	if notARepo {
		return Status{NotARepo: true}, nil
	}
	return status, nil
}

func cachedStatus(entry cachedEntry) Status {
	if entry.NotARepo {
		return Status{NotARepo: true}
	}
	return entry.Status
}

// cacheKey derives a filesystem-safe cache key from sessionID and dir
// (rather than the process PID, which is unique per invocation and would
// never hit) so a burst of refreshes within one session reuses the same
// cache entry. Hashing dir alongside sessionID handles a long-lived session
// that cd's between repositories.
func cacheKey(sessionID, dir string) string {
	sum := sha256.Sum256([]byte(sessionID + "|" + dir))
	return hex.EncodeToString(sum[:16])
}

func gitCacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("gitstatus: resolve cache dir: %w", err)
	}
	return filepath.Join(base, "statusline", "git"), nil
}

func headMTimeUnixNano(dir string) int64 {
	info, err := os.Stat(filepath.Join(dir, ".git", "HEAD"))
	if err != nil {
		return 0
	}
	return info.ModTime().UnixNano()
}

func runGitStatus(ctx context.Context, cfg config.GitConfig, dir string) (Status, bool, error) {
	timeout := time.Duration(cfg.TimeoutMS) * time.Millisecond
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(cctx, "git", "-C", dir, "status", "--porcelain=v2", "--branch", "--untracked-files=normal")
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0", "LC_ALL=C")
	out, err := cmd.Output()
	if err != nil {
		if _, ok := errors.AsType[*exec.ExitError](err); ok {
			// git exits non-zero for "not a git repository" and similar
			// usage errors; treat that as "no repo here" so the git
			// segment cleanly omits itself instead of surfacing an error.
			return Status{}, true, nil
		}
		return Status{}, false, fmt.Errorf("gitstatus: run git status: %w", err)
	}

	st, err := ParsePorcelainV2(out)
	if err != nil {
		return Status{}, false, err
	}
	return st, false, nil
}
