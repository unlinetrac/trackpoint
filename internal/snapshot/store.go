package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultStoreDir = ".trackpoint/snapshots"

// Store handles persistence of snapshots to disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at the given directory.
// If dir is empty, it defaults to .trackpoint/snapshots.
func NewStore(dir string) *Store {
	if dir == "" {
		dir = defaultStoreDir
	}
	return &Store{dir: dir}
}

// Save writes a snapshot to disk as <id>.json.
func (st *Store) Save(s *Snapshot) error {
	if err := os.MkdirAll(st.dir, 0755); err != nil {
		return fmt.Errorf("store: mkdir failed: %w", err)
	}
	data, err := s.Marshal()
	if err != nil {
		return err
	}
	path := filepath.Join(st.dir, s.ID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("store: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot by ID from disk.
func (st *Store) Load(id string) (*Snapshot, error) {
	path := filepath.Join(st.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("store: read failed for id %q: %w", id, err)
	}
	return Unmarshal(data)
}

// ListIDs returns the IDs of all stored snapshots.
func (st *Store) ListIDs() ([]string, error) {
	entries, err := os.ReadDir(st.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("store: readdir failed: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			ids = append(ids, e.Name()[:len(e.Name())-5])
		}
	}
	return ids, nil
}
