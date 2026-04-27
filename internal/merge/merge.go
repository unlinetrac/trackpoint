// Package merge provides functionality for combining two snapshots into one,
// applying a configurable conflict resolution strategy when the same key
// exists in both snapshots.
package merge

import (
	"errors"
	"fmt"
	"time"

	"github.com/trackpoint/internal/snapshot"
)

// Strategy controls how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyPreferBase keeps the value from the base snapshot on conflict.
	StrategyPreferBase Strategy = iota
	// StrategyPreferOther keeps the value from the other snapshot on conflict.
	StrategyPreferOther
	// StrategyError returns an error if any key exists in both snapshots.
	StrategyError
)

// ErrConflict is returned when StrategyError is used and a key collision
// is detected between the two snapshots.
var ErrConflict = errors.New("merge conflict: duplicate key")

// Options configures the behaviour of Run.
type Options struct {
	// Strategy determines how overlapping keys are handled.
	Strategy Strategy
	// Label is stored as the "merge" metadata field on the resulting snapshot.
	Label string
}

// Run merges base and other into a new snapshot according to opts.
// The resulting snapshot contains the union of both key sets; the supplied
// label (if any) is recorded under the "merge" metadata key alongside the
// source IDs so the lineage is traceable.
func Run(base, other *snapshot.Snapshot, opts Options) (*snapshot.Snapshot, error) {
	if base == nil {
		return nil, errors.New("merge: base snapshot must not be nil")
	}
	if other == nil {
		return nil, errors.New("merge: other snapshot must not be nil")
	}

	merged := make(map[string]string, len(base.State)+len(other.State))

	// Copy base state first.
	for k, v := range base.State {
		merged[k] = v
	}

	// Apply other state, respecting the conflict strategy.
	for k, v := range other.State {
		if existing, conflict := merged[k]; conflict {
			switch opts.Strategy {
			case StrategyPreferBase:
				// Keep the base value — nothing to do.
				_ = existing
			case StrategyPreferOther:
				merged[k] = v
			case StrategyError:
				return nil, fmt.Errorf("%w: %q (base=%q, other=%q)", ErrConflict, k, existing, v)
			}
		} else {
			merged[k] = v
		}
	}

	// Build metadata that records the merge provenance.
	meta := map[string]string{
		"merge.base":  base.ID,
		"merge.other": other.ID,
		"merge.time":  time.Now().UTC().Format(time.RFC3339),
	}
	if opts.Label != "" {
		meta["merge.label"] = opts.Label
	}

	snap := snapshot.New(merged, meta)
	return snap, nil
}
