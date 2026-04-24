package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Store manages reading and writing snapshots to disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at the given directory.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes a snapshot to disk as <id>.json.
func (s *Store) Save(snap *Snapshot) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("creating store dir: %w", err)
	}
	path := filepath.Join(s.dir, snap.ID+".json")
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing snapshot file: %w", err)
	}
	return nil
}

// Load reads a snapshot by ID from disk.
func (s *Store) Load(id string) (*Snapshot, error) {
	path := filepath.Join(s.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %w", err)
	}
	return Unmarshal(data)
}

// List returns all snapshot IDs in the store, sorted by filename.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading store dir: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".json" {
			ids = append(ids, name[:len(name)-5])
		}
	}
	sort.Strings(ids)
	return ids, nil
}
