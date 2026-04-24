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

// Store manages persistence of snapshots on the local filesystem.
// Each snapshot is stored as a JSON file named by its ID within a directory.
type Store struct {
	dir string
}

// NewStore creates a Store that reads and writes snapshots under the given
// directory. The directory is created lazily on first write.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save persists a snapshot to disk. The file is named <id>.json.
// The storage directory is created if it does not already exist.
func (s *Store) Save(snap *Snapshot) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("snapshot store: create dir: %w", err)
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot store: marshal: %w", err)
	}

	path := filepath.Join(s.dir, snap.ID+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot store: write file: %w", err)
	}
	return nil
}

// Load retrieves a snapshot by ID. Returns an error wrapping os.ErrNotExist
// when no snapshot with that ID is found.
func (s *Store) Load(id string) (*Snapshot, error) {
	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("snapshot store: id %q not found: %w", id, os.ErrNotExist)
		}
		return nil, fmt.Errorf("snapshot store: read file: %w", err)
	}

	snap, err := Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: unmarshal %q: %w", id, err)
	}
	return snap, nil
}

// List returns all snapshot IDs stored in the directory, sorted lexicographically.
// Returns an empty slice (and no error) when the directory does not exist yet.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("snapshot store: read dir: %w", err)
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
