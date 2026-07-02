package cache

import (
	"os"
	"path/filepath"
	"testing"
)

type entry struct {
	Branch string `json:"branch"`
	Count  int    `json:"count"`
}

func TestStoreAtomicAndLoadRoundTrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	want := entry{Branch: "main", Count: 3}
	if err := StoreAtomic(dir, "gitstatus", want); err != nil {
		t.Fatalf("StoreAtomic() error = %v", err)
	}

	got, ok := Load[entry](dir, "gitstatus")
	if !ok {
		t.Fatal("Load() ok = false, want true")
	}
	if got != want {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestStoreAtomicLeavesNoTempFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if err := StoreAtomic(dir, "gitstatus", entry{Branch: "main"}); err != nil {
		t.Fatalf("StoreAtomic() error = %v", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}
	if len(files) != 1 || files[0].Name() != "gitstatus" {
		t.Errorf("dir contents = %v, want exactly [gitstatus]", files)
	}
}

func TestStoreAtomicCreatesDir(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "nested", "cache", "dir")

	if err := StoreAtomic(dir, "key", entry{Branch: "main"}); err != nil {
		t.Fatalf("StoreAtomic() error = %v", err)
	}
	if _, ok := Load[entry](dir, "key"); !ok {
		t.Error("Load() ok = false after StoreAtomic created the dir, want true")
	}
}

func TestLoad_missingFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if _, ok := Load[entry](dir, "does-not-exist"); ok {
		t.Error("Load() ok = true for missing file, want false")
	}
}

func TestLoad_corruptJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad"), []byte("not json"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, ok := Load[entry](dir, "bad"); ok {
		t.Error("Load() ok = true for corrupt JSON, want false")
	}
}

func TestStoreAtomic_mkdirFailsWhenParentIsAFile(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	blocker := filepath.Join(base, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	dir := filepath.Join(blocker, "cache") // blocker is a file, not a dir

	if err := StoreAtomic(dir, "key", entry{Branch: "main"}); err == nil {
		t.Error("StoreAtomic() error = nil, want error when a path component is a file")
	}
}

func TestStoreAtomic_marshalFailsForUnsupportedType(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if err := StoreAtomic(dir, "key", make(chan int)); err == nil {
		t.Error("StoreAtomic() error = nil, want error for an unmarshalable type")
	}
}

func TestStoreAtomic_createTempFailsWhenDirNotWritable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission checks don't apply")
	}
	t.Parallel()
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) }) // let t.TempDir() clean up afterward

	if err := StoreAtomic(dir, "key", entry{Branch: "main"}); err == nil {
		t.Error("StoreAtomic() error = nil, want error when dir isn't writable")
	}
}

func TestStoreAtomic_renameFailsWhenDestinationIsANonEmptyDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	destDir := filepath.Join(dir, "key")
	if err := os.MkdirAll(filepath.Join(destDir, "inner"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if err := StoreAtomic(dir, "key", entry{Branch: "main"}); err == nil {
		t.Error("StoreAtomic() error = nil, want error when the destination is a non-empty directory")
	}
}

func TestStoreAtomicOverwritesExisting(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	if err := StoreAtomic(dir, "key", entry{Branch: "main", Count: 1}); err != nil {
		t.Fatalf("StoreAtomic() error = %v", err)
	}
	if err := StoreAtomic(dir, "key", entry{Branch: "develop", Count: 2}); err != nil {
		t.Fatalf("StoreAtomic() error = %v", err)
	}

	got, ok := Load[entry](dir, "key")
	if !ok {
		t.Fatal("Load() ok = false, want true")
	}
	want := entry{Branch: "develop", Count: 2}
	if got != want {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}
