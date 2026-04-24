package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Store manages persistence of snapshots on disk.
// Each snapshot is stored as a JSON file named by its ID.
type Store struct {
	dir string
}

// NewStore creates a Store that reads and writes snapshots under dir.
// The directory is created lazily on first write.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes s to disk as <dir>/<s.ID>.json, creating the directory if needed.
func (s *Store) Save(snap *Snapshot) error {
	if snap == nil {
		return errors.New("snapshot: cannot save nil snapshot")
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: create store dir: %w", err)
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := s.snapshotPath(snap.ID)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads the snapshot with the given id from disk.
// Returns an error wrapping os.ErrNotExist if the snapshot is not found.
func (s *Store) Load(id string) (*Snapshot, error) {
	path := s.snapshotPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("snapshot %q not found: %w", id, os.ErrNotExist)
		}
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	snap, err := Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal %s: %w", path, err)
	}
	return snap, nil
}

// List returns all snapshot IDs present in the store, sorted lexicographically.
// Returns an empty slice (and no error) when the store directory does not exist.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("snapshot: list store dir: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".json") {
			ids = append(ids, strings.TrimSuffix(name, ".json"))
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// snapshotPath returns the full file path for a snapshot with the given id.
func (s *Store) snapshotPath(id string) string {
	return filepath.Join(s.dir, id+".json")
}
