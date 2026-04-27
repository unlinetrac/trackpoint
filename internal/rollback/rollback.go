// Package rollback provides utilities for identifying a safe snapshot
// to roll back to based on a baseline or a named alias.
package rollback

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
)

// SnapshotLoader is satisfied by snapshot.Store.
type SnapshotLoader interface {
	Load(id string) (*snapshot.Snapshot, error)
	List() ([]string, error)
}

// Result describes the outcome of a rollback target search.
type Result struct {
	// Target is the snapshot recommended for rollback.
	Target *snapshot.Snapshot
	// Changes is the diff between Target and the current snapshot.
	Changes []diff.Change
}

// ErrNoCandidate is returned when no suitable rollback target is found.
var ErrNoCandidate = errors.New("rollback: no suitable candidate found")

// Find searches backwards through stored snapshots (excluding currentID)
// and returns the most recent snapshot whose diff against current contains
// no changes, or the immediately preceding snapshot when noDiff is false.
func Find(store SnapshotLoader, currentID string, noDiff bool) (*Result, error) {
	ids, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("rollback: list snapshots: %w", err)
	}

	current, err := store.Load(currentID)
	if err != nil {
		return nil, fmt.Errorf("rollback: load current snapshot %q: %w", currentID, err)
	}

	// Walk in reverse (newest first), skip currentID.
	for i := len(ids) - 1; i >= 0; i-- {
		if ids[i] == currentID {
			continue
		}
		candidate, err := store.Load(ids[i])
		if err != nil {
			continue
		}
		result, err := diff.Compare(candidate, current)
		if err != nil {
			continue
		}
		if noDiff && len(result.Changes) > 0 {
			continue
		}
		return &Result{Target: candidate, Changes: result.Changes}, nil
	}
	return nil, ErrNoCandidate
}
