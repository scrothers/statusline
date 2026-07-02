package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Load reads and JSON-decodes the cache entry named key from dir into a T.
// It reports ok=false — never an error — for any failure (missing file,
// unreadable, corrupt JSON), since a cache miss and a cache read failure are
// handled identically by every caller: fall through to recomputing.
func Load[T any](dir, key string) (entry T, ok bool) {
	data, err := os.ReadFile(filepath.Join(dir, key))
	if err != nil {
		return entry, false
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return entry, false
	}
	return entry, true
}

// StoreAtomic JSON-encodes entry and writes it to dir/key, via a temp file
// in the same directory renamed into place, so a process killed mid-write
// (Claude Code SIGTERMs a superseded in-flight statusline run) never leaves
// a partially-written cache file for the next invocation to trip over.
func StoreAtomic[T any](dir, key string, entry T) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cache: create dir %s: %w", dir, err)
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cache: marshal %s: %w", key, err)
	}

	tmp, err := os.CreateTemp(dir, key+".tmp-*")
	if err != nil {
		return fmt.Errorf("cache: create temp file for %s: %w", key, err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }() // no-op once the rename below succeeds

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("cache: write temp file for %s: %w", key, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("cache: close temp file for %s: %w", key, err)
	}
	if err := os.Rename(tmpPath, filepath.Join(dir, key)); err != nil {
		return fmt.Errorf("cache: rename into place for %s: %w", key, err)
	}
	return nil
}
