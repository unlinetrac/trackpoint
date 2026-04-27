// Package compare provides multi-snapshot comparison across a time range,
// summarising net changes between the earliest and latest snapshot in a set.
package compare

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
)

// Result holds the outcome of a multi-snapshot range comparison.
type Result struct {
	// FromID is the ID of the earliest snapshot used as the baseline.
	FromID string
	// ToID is the ID of the latest snapshot compared against the baseline.
	ToID string
	// Diff contains the net change between From and To.
	Diff diff.Result
	// SnapshotCount is the total number of snapshots inspected.
	SnapshotCount int
}

// HasChanges reports whether any net changes exist between the two endpoints.
func (r Result) HasChanges() bool {
	return len(r.Diff.Added)+len(r.Diff.Removed)+len(r.Diff.Modified) > 0
}

// Run compares the first and last snapshots in the provided slice.
// Snapshots must be ordered oldest-first. At least two snapshots are required.
func Run(snapshots []*snapshot.Snapshot) (Result, error) {
	if len(snapshots) < 2 {
		return Result{}, errors.New("compare: at least two snapshots are required")
	}

	from := snapshots[0]
	to := snapshots[len(snapshots)-1]

	d, err := diff.Compare(from, to)
	if err != nil {
		return Result{}, fmt.Errorf("compare: diff failed: %w", err)
	}

	return Result{
		FromID:        from.ID,
		ToID:          to.ID,
		Diff:          d,
		SnapshotCount: len(snapshots),
	}, nil
}
