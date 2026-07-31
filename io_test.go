package keystone

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsurePathWithinRoot(t *testing.T) {
	const root = "/var/app/data"

	tests := []struct {
		name    string
		path    string
		root    string
		wantErr bool
	}{
		{
			name: "path equals root",
			path: "/var/app/data",
			root: root,
		},
		{
			name: "root with trailing slash still matches",
			path: "/var/app/data/",
			root: root,
		},
		{
			name: "simple subdirectory",
			path: "/var/app/data/sub/file.txt",
			root: root,
		},
		{
			name: "dot-dot that stays inside root after cleaning",
			path: "/var/app/data/sub/../other/file.txt",
			root: root,
		},
		{
			name:    "dot-dot that escapes root entirely",
			path:    "/var/app/data/../secret",
			root:    root,
			wantErr: true,
		},
		{
			name:    "path is a parent of root",
			path:    "/var/app",
			root:    root,
			wantErr: true,
		},
		{
			name:    "unrelated absolute path",
			path:    "/etc/passwd",
			root:    root,
			wantErr: true,
		},
		{
			name:    "sibling directory sharing root as a string prefix",
			path:    "/var/app/data-evil/file.txt",
			root:    root,
			wantErr: true,
		},
		{
			name: "deeply nested traversal that nets out inside root",
			path: "/var/app/data/a/b/../../a/file.txt",
			root: root,
		},
		{
			name:    "traversal several levels above root",
			path:    "/var/app/data/../../../etc/passwd",
			root:    root,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EnsurePathWithinRoot(tt.path, tt.root)
			if tt.wantErr && err == nil {
				t.Fatalf("EnsurePathWithinRoot(%q, %q) = nil, want error", tt.path, tt.root)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("EnsurePathWithinRoot(%q, %q) = %v, want nil", tt.path, tt.root, err)
			}
			if tt.wantErr && !errors.Is(err, ErrPathTraversal) {
				t.Fatalf("EnsurePathWithinRoot(%q, %q) error = %v, want it to wrap ErrPathTraversal", tt.path, tt.root, err)
			}
		})
	}
}

// TestEnsurePathWithinRoot_RelativePaths confirms relative inputs are
// resolved against the current working directory, by actually changing into
// a temp directory rather than relying on string comparisons.
func TestEnsurePathWithinRoot_RelativePaths(t *testing.T) {
	tmp := t.TempDir()
	dataDir := filepath.Join(tmp, "data")
	if err := os.Mkdir(dataDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if _, err := EnsurePathWithinRoot("./data/file.txt", "./data"); err != nil {
		t.Errorf("relative subpath should be allowed, got: %v", err)
	}

	if _, err := EnsurePathWithinRoot("./data/../outside.txt", "./data"); err == nil {
		t.Errorf("relative escape should be rejected, got nil")
	}
}

// TestEnsurePathWithinRoot_SymlinkLimitation documents a known limitation:
// EnsurePathWithinRoot operates purely on lexical paths and does not follow
// symlinks. A symlink inside root that points outside it will NOT be caught
// here — callers handling untrusted symlinks must resolve them first with
// filepath.EvalSymlinks and check the resolved target instead.
func TestEnsurePathWithinRoot_SymlinkLimitation(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "root")
	outside := filepath.Join(tmp, "outside")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	link := filepath.Join(root, "escape")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks not supported in this environment: %v", err)
	}

	// Lexically, "root/escape" is inside root, so this passes even though
	// the symlink target is not. That's the documented limitation.
	if _, err := EnsurePathWithinRoot(link, root); err != nil {
		t.Fatalf("expected lexical check to pass (documenting the limitation), got: %v", err)
	}

	// The fix: resolve symlinks before validating.
	resolved, err := filepath.EvalSymlinks(link)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	if _, err := EnsurePathWithinRoot(resolved, root); err == nil {
		t.Fatalf("expected resolved symlink target outside root to be rejected")
	}
}
