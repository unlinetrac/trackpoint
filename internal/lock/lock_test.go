package lock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/trackpoint/internal/lock"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "lock-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLock_AcquireAndRelease(t *testing.T) {
	s := lock.NewStore(tempDir(t))
	const id = "snap-abc"

	if err := s.Lock(id); err != nil {
		t.Fatalf("Lock: unexpected error: %v", err)
	}
	if !s.IsLocked(id) {
		t.Fatal("expected snapshot to be locked")
	}
	if err := s.Unlock(id); err != nil {
		t.Fatalf("Unlock: unexpected error: %v", err)
	}
	if s.IsLocked(id) {
		t.Fatal("expected snapshot to be unlocked after Unlock")
	}
}

func TestLock_DoubleLock_ReturnsErrLocked(t *testing.T) {
	s := lock.NewStore(tempDir(t))
	const id = "snap-xyz"

	if err := s.Lock(id); err != nil {
		t.Fatalf("first Lock: %v", err)
	}
	defer s.Unlock(id) //nolint:errcheck

	err := s.Lock(id)
	if err == nil {
		t.Fatal("expected ErrLocked, got nil")
	}
	if err != lock.ErrLocked {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestUnlock_NotFound_NoError(t *testing.T) {
	s := lock.NewStore(tempDir(t))
	if err := s.Unlock("nonexistent"); err != nil {
		t.Fatalf("Unlock on missing lock should not error: %v", err)
	}
}

func TestIsLocked_CreatesLockFile(t *testing.T) {
	dir := tempDir(t)
	s := lock.NewStore(dir)
	const id = "snap-check"

	if s.IsLocked(id) {
		t.Fatal("expected not locked before Lock")
	}
	_ = s.Lock(id)

	if _, err := os.Stat(filepath.Join(dir, id+".lock")); err != nil {
		t.Fatalf("expected lock file to exist: %v", err)
	}
}
