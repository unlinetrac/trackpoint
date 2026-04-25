package baseline

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/snapshot"
)

// ErrNoBaseline is returned when no baseline has been set.
var ErrNoBaseline = errors.New("no baseline set")

// Store manages a single "baseline" snapshot that diffs are measured against.
type Store struct {
	store snapshotStore
	baselineKey string
}

type snapshotStore interface {
	Save(s *snapshot.Snapshot) error
	Load(id string) (*snapshot.Snapshot, error)
}

// New creates a new baseline Store backed by the given snapshot store.
func New(store snapshotStore) *Store {
	return &Store{
		store:       store,
		baselineKey: "__baseline__",
	}
}

// Set persists the given snapshot as the current baseline.
func (b *Store) Set(s *snapshot.Snapshot) error {
	if s == nil {
		return errors.New("snapshot must not be nil")
	}
	// Clone with a stable baseline ID so it doesn't collide with real snapshots.
	copy := *s
	copy.ID = b.baselineKey
	return b.store.Save(&copy)
}

// Get retrieves the current baseline snapshot.
func (b *Store) Get() (*snapshot.Snapshot, error) {
	s, err := b.store.Load(b.baselineKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoBaseline, err)
	}
	return s, nil
}

// Clear removes the baseline by overwriting with an empty snapshot.
func (b *Store) Clear() error {
	empty := &snapshot.Snapshot{
		ID:   b.baselineKey,
		Data: map[string]string{},
	}
	return b.store.Save(empty)
}

// IsSet returns true when a non-empty baseline exists.
func (b *Store) IsSet() bool {
	s, err := b.Get()
	if err != nil {
		return false
	}
	return len(s.Data) > 0
}
