// Package prune provides utilities for removing old snapshots from the store
// based on retention policies such as keeping only the N most recent snapshots.
package prune

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/snapshot"
)

// Store is the interface required by the pruner to list and delete snapshots.
type Store interface {
	List() ([]string, error)
	Delete(id string) error
}

// Options configures how pruning is performed.
type Options struct {
	// KeepLast is the number of most-recent snapshots to retain.
	// Must be >= 1.
	KeepLast int
	// DryRun reports what would be deleted without actually deleting.
	DryRun bool
}

// Result holds the outcome of a prune operation.
type Result struct {
	Deleted []string
	Skipped []string
}

// Run applies the retention policy described by opts against the given store.
// Snapshots are assumed to be returned by List in ascending chronological order
// (oldest first), which matches the sorted-ID contract of snapshot.Store.
func Run(store Store, opts Options) (*Result, error) {
	if opts.KeepLast < 1 {
		return nil, errors.New("prune: KeepLast must be at least 1")
	}

	ids, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("prune: listing snapshots: %w", err)
	}

	result := &Result{}

	if len(ids) <= opts.KeepLast {
		result.Skipped = append(result.Skipped, ids...)
		return result, nil
	}

	cutoff := len(ids) - opts.KeepLast
	toDelete := ids[:cutoff]
	toKeep := ids[cutoff:]

	result.Skipped = toKeep

	for _, id := range toDelete {
		if opts.DryRun {
			result.Deleted = append(result.Deleted, id)
			continue
		}
		if err := store.Delete(id); err != nil {
			return result, fmt.Errorf("prune: deleting snapshot %q: %w", id, err)
		}
		result.Deleted = append(result.Deleted, id)
	}

	return result, nil
}

// ensure snapshot.Store satisfies our interface at compile time.
var _ Store = (*snapshot.Store)(nil)
