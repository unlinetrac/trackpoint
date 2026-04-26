// Package alias provides named aliases for snapshot IDs, allowing
// human-readable references like "prod-v1" instead of raw SHA IDs.
package alias

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var validAlias = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ErrNotFound is returned when an alias does not exist.
var ErrNotFound = errors.New("alias not found")

// Store persists alias-to-snapshot-ID mappings on disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Set associates name with snapshotID, overwriting any existing mapping.
func (s *Store) Set(name, snapshotID string) error {
	if !validAlias.MatchString(name) {
		return fmt.Errorf("invalid alias %q: must match [a-zA-Z0-9_-]+", name)
	}
	if snapshotID == "" {
		return errors.New("snapshot ID must not be empty")
	}
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("create alias dir: %w", err)
	}
	data, err := json.Marshal(snapshotID)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(name), data, 0644)
}

// Get returns the snapshot ID for the given alias name.
func (s *Store) Get(name string) (string, error) {
	data, err := os.ReadFile(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return "", err
	}
	var id string
	if err := json.Unmarshal(data, &id); err != nil {
		return "", fmt.Errorf("corrupt alias file %q: %w", name, err)
	}
	return id, nil
}

// Delete removes the alias. Returns nil if it does not exist.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// List returns all defined alias names in no guaranteed order.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name)
}
