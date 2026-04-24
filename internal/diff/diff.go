package diff

import (
	"fmt"
	"sort"

	"github.com/trackpoint/internal/snapshot"
)

// ChangeType represents the kind of change detected.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single key-level difference between two snapshots.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal interface{}
	NewVal interface{}
}

// Result holds the full diff between two snapshots.
type Result struct {
	FromID  string
	ToID    string
	Changes []Change
}

// HasChanges returns true if any changes were detected.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a human-readable summary of the diff.
func (r *Result) Summary() string {
	if !r.HasChanges() {
		return fmt.Sprintf("No changes between %s and %s", r.FromID, r.ToID)
	}
	return fmt.Sprintf("%d change(s) between %s and %s", len(r.Changes), r.FromID, r.ToID)
}

// Compare computes the diff between two snapshots.
func Compare(from, to *snapshot.Snapshot) *Result {
	result := &Result{
		FromID: from.ID,
		ToID:   to.ID,
	}

	for key, newVal := range to.State {
		oldVal, exists := from.State[key]
		if !exists {
			result.Changes = append(result.Changes, Change{Key: key, Type: Added, NewVal: newVal})
		} else if fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal) {
			result.Changes = append(result.Changes, Change{Key: key, Type: Modified, OldVal: oldVal, NewVal: newVal})
		}
	}

	for key, oldVal := range from.State {
		if _, exists := to.State[key]; !exists {
			result.Changes = append(result.Changes, Change{Key: key, Type: Removed, OldVal: oldVal})
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Key < result.Changes[j].Key
	})

	return result
}
