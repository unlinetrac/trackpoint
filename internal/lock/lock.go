// Package lock provides snapshot locking to prevent concurrent modifications.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrLocked is returned when a snapshot is already locked.
var ErrLocked = errors.New("snapshot is locked")

// Store manages lock files for snapshots.
type Store struct {
	dir string
}

// NewStore creates a new lock Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Lock acquires a lock for the given snapshot ID.
// Returns ErrLocked if the snapshot is already locked.
func (s *Store) Lock(snapshotID string) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("lock: create dir: %w", err)
	}
	path := s.lockPath(snapshotID)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return ErrLocked
		}
		return fmt.Errorf("lock: open file: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d", time.Now().UnixNano())
	return err
}

// Unlock releases the lock for the given snapshot ID.
// Returns nil if no lock exists.
func (s *Store) Unlock(snapshotID string) error {
	path := s.lockPath(snapshotID)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("lock: remove: %w", err)
	}
	return nil
}

// IsLocked reports whether the given snapshot ID is currently locked.
func (s *Store) IsLocked(snapshotID string) bool {
	_, err := os.Stat(s.lockPath(snapshotID))
	return err == nil
}

func (s *Store) lockPath(snapshotID string) string {
	return filepath.Join(s.dir, snapshotID+".lock")
}
