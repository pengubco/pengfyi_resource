package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeCache lays out a minimal containerd root in dir, as the image cache
// tarball would, and returns the paths of the entries that link must handle.
func writeCache(t *testing.T, dir string) {
	t.Helper()
	files := map[string]string{
		filepath.Join(metadataDir, "meta.db"):                         "meta",
		filepath.Join(contentDir, "aaaa"):                             "blob-a",
		filepath.Join(contentDir, "bbbb"):                             "blob-b",
		filepath.Join(snapshotterDir, snapshotterDB):                  "snapshot-meta",
		filepath.Join(snapshotterDir, snapshotsDir, "1", "fs", "bin"): "layer-1",
		filepath.Join(snapshotterDir, snapshotsDir, "2", "fs", "etc"): "layer-2",
	}
	for rel, content := range files {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// containerd creates the snapshotter db with restrictive permissions;
	// make sure they survive the copy.
	if err := os.Chmod(filepath.Join(dir, snapshotterDir, snapshotterDB), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestLink(t *testing.T) {
	cache := t.TempDir()
	root := t.TempDir()
	writeCache(t, cache)

	if err := link(root, cache); err != nil {
		t.Fatalf("link: %v", err)
	}

	// Databases are copies, not links, with permissions preserved.
	for _, rel := range []string{
		filepath.Join(metadataDir, "meta.db"),
		filepath.Join(snapshotterDir, snapshotterDB),
	} {
		info, err := os.Lstat(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			t.Errorf("%s: want a regular file, got a symlink", rel)
		}
		want, _ := os.Stat(filepath.Join(cache, rel))
		if info.Mode().Perm() != want.Mode().Perm() {
			t.Errorf("%s: mode %v, want %v", rel, info.Mode().Perm(), want.Mode().Perm())
		}
	}

	// Blobs and snapshots are symlinks into the cache.
	for _, rel := range []string{
		filepath.Join(contentDir, "aaaa"),
		filepath.Join(contentDir, "bbbb"),
		filepath.Join(snapshotterDir, snapshotsDir, "1"),
		filepath.Join(snapshotterDir, snapshotsDir, "2"),
	} {
		got, err := os.Readlink(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
		if want := filepath.Join(cache, rel); got != want {
			t.Errorf("%s -> %s, want %s", rel, got, want)
		}
	}

	// Reading through a snapshot link reaches the cached layer.
	data, err := os.ReadFile(filepath.Join(root, snapshotterDir, snapshotsDir, "2", "fs", "etc"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "layer-2" {
		t.Errorf("read %q through snapshot link, want layer-2", data)
	}
}

func TestLinkEmptiesRootOnFailure(t *testing.T) {
	cache := t.TempDir()
	root := t.TempDir()
	writeCache(t, cache)
	// Remove the snapshotter db, which is handled last, so the earlier steps
	// have already written into root when the failure happens.
	if err := os.Remove(filepath.Join(cache, snapshotterDir, snapshotterDB)); err != nil {
		t.Fatal(err)
	}

	if err := link(root, cache); err == nil {
		t.Fatal("link succeeded with an incomplete cache")
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("root not emptied after failure: %d entries left", len(entries))
	}
}

func TestLinkRefusesToOverwrite(t *testing.T) {
	cache := t.TempDir()
	root := t.TempDir()
	writeCache(t, cache)
	// A root with its own meta.db belongs to a containerd that has already
	// run; populating or emptying it would discard its state.
	if err := os.MkdirAll(filepath.Join(root, metadataDir), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, metadataDir, "meta.db"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := link(root, cache); err == nil {
		t.Fatal("link populated a non-empty root")
	}

	data, err := os.ReadFile(filepath.Join(root, metadataDir, "meta.db"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "existing" {
		t.Errorf("existing meta.db was modified: %q", data)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 1 {
		t.Errorf("root was modified: %d entries", len(entries))
	}
}
