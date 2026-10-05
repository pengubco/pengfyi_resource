// image-cache-link prepares an empty containerd root directory so that
// containerd finds the images of a read-only image cache in it.
//
// containerd keeps four things under its root that matter for images:
//
//	io.containerd.metadata.v1.bolt/meta.db                 images, references, leases
//	io.containerd.content.v1.content/blobs/sha256/<digest> manifests, configs, layers
//	io.containerd.snapshotter.v1.overlayfs/metadata.db     snapshot tree
//	io.containerd.snapshotter.v1.overlayfs/snapshots/<id>  unpacked layers
//
// The two databases are opened read-write, so they are copied to the new root.
// Blobs and snapshots are immutable once written, so each entry becomes a
// symlink into the cache. containerd follows the symlinks and never knows the
// bytes live on another volume.
//
// On any error the new root is emptied again, so a broken cache leaves a
// working, if cold, containerd behind.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	metadataDir    = "io.containerd.metadata.v1.bolt"
	contentDir     = "io.containerd.content.v1.content/blobs/sha256"
	snapshotterDir = "io.containerd.snapshotter.v1.overlayfs"
	snapshotsDir   = "snapshots"
	snapshotterDB  = "metadata.db"
)

func main() {
	root := flag.String("root", "/var/lib/containerd", "containerd root directory to prepare (must be empty)")
	cache := flag.String("cache", "/var/lib/image-cache", "image cache directory with the layout of a containerd root")
	flag.Parse()

	if err := link(*root, *cache); err != nil {
		fmt.Fprintf(os.Stderr, "image-cache-link: %v\n", err)
		os.Exit(1)
	}
}

// link populates root from cache, and empties root again if that fails.
// A root that already has content is left alone: it belongs to a containerd
// that has run before, and nothing in it should be discarded here.
func link(root, cache string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return fmt.Errorf("%s is not empty, refusing to populate it", root)
	}

	err = populate(root, cache)
	if err == nil {
		return nil
	}
	if cleanupErr := emptyDir(root); cleanupErr != nil {
		return errors.Join(err, fmt.Errorf("restoring %s to an empty directory: %w", root, cleanupErr))
	}
	return err
}

func populate(root, cache string) error {
	if err := copyDir(filepath.Join(cache, metadataDir), filepath.Join(root, metadataDir)); err != nil {
		return fmt.Errorf("copying metadata: %w", err)
	}
	if err := symlinkEntries(filepath.Join(cache, contentDir), filepath.Join(root, contentDir)); err != nil {
		return fmt.Errorf("linking content blobs: %w", err)
	}
	if err := symlinkEntries(filepath.Join(cache, snapshotterDir, snapshotsDir), filepath.Join(root, snapshotterDir, snapshotsDir)); err != nil {
		return fmt.Errorf("linking snapshots: %w", err)
	}
	if err := copyFile(filepath.Join(cache, snapshotterDir, snapshotterDB), filepath.Join(root, snapshotterDir, snapshotterDB)); err != nil {
		return fmt.Errorf("copying snapshotter metadata: %w", err)
	}
	return nil
}

// symlinkEntries creates dst with the mode of src, and one symlink in dst for
// each entry of src, pointing back at the entry in src.
func symlinkEntries(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		target := filepath.Join(src, entry.Name())
		if err := os.Symlink(target, filepath.Join(dst, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

// copyDir copies the regular files and directories under src to dst,
// preserving permission bits. Extended attributes, including SELinux labels,
// are intentionally not copied: dst takes the labels of its own filesystem.
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		switch {
		case d.IsDir():
			return os.MkdirAll(target, info.Mode().Perm())
		case d.Type().IsRegular():
			return copyFile(path, target)
		default:
			return fmt.Errorf("%s: unsupported file type %s", path, d.Type())
		}
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func emptyDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var errs []error
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
