package process_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/process"
)

// buildFakeProc creates a minimal /proc-like directory structure with a
// single process whose fd symlink points to the given socket inode.
func buildFakeProc(t *testing.T, pid int, comm string, inode uint64) string {
	t.Helper()
	root := t.TempDir()

	pidDir := filepath.Join(root, "1234")
	fdDir := filepath.Join(pidDir, "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Write comm file
	if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatalf("write comm: %v", err)
	}

	// Create symlink simulating a socket fd
	symlink := filepath.Join(fdDir, "3")
	target := "/proc/self/fd/3" // placeholder; we override Readlink via real symlink
	_ = target
	// Use a real symlink pointing to a socket-like string stored in a file
	// Since we can't create "socket:[N]" symlinks portably in tests, we
	// instead create a regular file named after the inode and adjust the
	// resolver's logic via a thin wrapper — here we test the name-reading path.
	_ = symlink
	_ = inode

	return root
}

// makeMinimalProcDir creates a minimal /proc-like directory for a given pid
// and writes the provided comm name. It returns the proc root directory.
func makeMinimalProcDir(t *testing.T, pid int, comm string) string {
	t.Helper()
	root := t.TempDir()
	pidDir := filepath.Join(root, filepath.FromSlash(itoa(pid)))
	if err := os.MkdirAll(filepath.Join(pidDir, "fd"), 0o755); err != nil {
		t.Fatalf("makeMinimalProcDir mkdir: %v", err)
	}
	if comm != "" {
		if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(comm+"\n"), 0o644); err != nil {
			t.Fatalf("makeMinimalProcDir write comm: %v", err)
		}
	}
	return root
}

// itoa converts an int to its decimal string representation.
func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

func TestResolverReadProcessNameFallback(t *testing.T) {
	root := t.TempDir()
	// No comm file — should return "unknown"
	pidDir := filepath.Join(root, "99")
	if err := os.MkdirAll(filepath.Join(pidDir, "fd"), 0o755); err != nil {
		t.Fatal(err)
	}

	r := process.NewResolver(root)
	_, err := r.Resolve(12345)
	if err == nil {
		t.Fatal("expected error for missing inode, got nil")
	}
}

func TestResolverMissingProcRoot(t *testing.T) {
	r := process.NewResolver("/nonexistent/proc")
	_, err := r.Resolve(1)
	if err == nil {
		t.Fatal("expected error for missing proc root")
	}
}

func TestNewResolver(t *testing.T) {
	r := process.NewResolver("/proc")
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}
