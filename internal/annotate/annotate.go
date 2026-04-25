// Package annotate provides functionality for attaching free-form notes
// to snapshots, stored alongside tag and baseline metadata.
package annotate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Annotation holds a text note associated with a snapshot ID.
type Annotation struct {
	SnapshotID string `json:"snapshot_id"`
	Note       string `json:"note"`
}

// Store persists annotations to a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a new annotation Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Set writes an annotation note for the given snapshot ID.
// Passing an empty note deletes any existing annotation.
func (s *Store) Set(snapshotID, note string) error {
	if snapshotID == "" {
		return errors.New("annotate: snapshot ID must not be empty")
	}
	if note == "" {
		return s.Delete(snapshotID)
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("annotate: create dir: %w", err)
	}
	a := Annotation{SnapshotID: snapshotID, Note: note}
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	return os.WriteFile(s.path(snapshotID), data, 0o644)
}

// Get retrieves the annotation for the given snapshot ID.
// Returns an empty Annotation (no error) if none exists.
func (s *Store) Get(snapshotID string) (Annotation, error) {
	data, err := os.ReadFile(s.path(snapshotID))
	if errors.Is(err, os.ErrNotExist) {
		return Annotation{SnapshotID: snapshotID}, nil
	}
	if err != nil {
		return Annotation{}, fmt.Errorf("annotate: read: %w", err)
	}
	var a Annotation
	if err := json.Unmarshal(data, &a); err != nil {
		return Annotation{}, fmt.Errorf("annotate: unmarshal: %w", err)
	}
	return a, nil
}

// Delete removes the annotation for the given snapshot ID.
// Returns nil if no annotation exists.
func (s *Store) Delete(snapshotID string) error {
	err := os.Remove(s.path(snapshotID))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *Store) path(snapshotID string) string {
	return filepath.Join(s.dir, snapshotID+".annotation.json")
}
