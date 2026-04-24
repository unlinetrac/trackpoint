package tag

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const tagsFile = "tags.json"

// TagStore persists snapshot tags to disk alongside snapshots.
type TagStore struct {
	dir string
}

// NewTagStore creates a TagStore rooted at dir.
func NewTagStore(dir string) *TagStore {
	return &TagStore{dir: dir}
}

// Save writes tags for the given snapshot ID to disk.
func (s *TagStore) Save(id string, tags Tags) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(s.dir, id+"."+tagsFile)
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads tags for the given snapshot ID from disk.
// Returns an empty Tags map if no tags file exists.
func (s *TagStore) Load(id string) (Tags, error) {
	path := filepath.Join(s.dir, id+"."+tagsFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Tags{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tags Tags
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// Delete removes the tags file for the given snapshot ID.
func (s *TagStore) Delete(id string) error {
	path := filepath.Join(s.dir, id+"."+tagsFile)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
