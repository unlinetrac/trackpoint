package history

import (
	"fmt"
	"sort"

	"github.com/user/trackpoint/internal/snapshot"
)

// Entry represents a single entry in the snapshot history timeline.
type Entry struct {
	Index    int
	Snapshot *snapshot.Snapshot
}

// Timeline holds an ordered sequence of snapshots for comparison.
type Timeline struct {
	entries []Entry
}

// NewTimeline builds a Timeline from a store and a list of snapshot IDs.
// IDs are resolved in the order provided; duplicates are ignored.
func NewTimeline(store *snapshot.Store, ids []string) (*Timeline, error) {
	seen := make(map[string]bool)
	var entries []Entry

	for i, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true

		snap, err := store.Load(id)
		if err != nil {
			return nil, fmt.Errorf("history: loading snapshot %q: %w", id, err)
		}
		entries = append(entries, Entry{Index: i, Snapshot: snap})
	}

	return &Timeline{entries: entries}, nil
}

// Entries returns the ordered list of timeline entries.
func (t *Timeline) Entries() []Entry {
	return t.entries
}

// Len returns the number of snapshots in the timeline.
func (t *Timeline) Len() int {
	return len(t.entries)
}

// Pairs returns consecutive (before, after) snapshot pairs for diffing.
func (t *Timeline) Pairs() [][2]*snapshot.Snapshot {
	if len(t.entries) < 2 {
		return nil
	}
	var pairs [][2]*snapshot.Snapshot
	for i := 0; i < len(t.entries)-1; i++ {
		pairs = append(pairs, [2]*snapshot.Snapshot{
			t.entries[i].Snapshot,
			t.entries[i+1].Snapshot,
		})
	}
	return pairs
}

// SortedByTime returns a new Timeline with entries sorted by snapshot timestamp.
func (t *Timeline) SortedByTime() *Timeline {
	copied := make([]Entry, len(t.entries))
	copy(copied, t.entries)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Snapshot.CreatedAt.Before(copied[j].Snapshot.CreatedAt)
	})
	return &Timeline{entries: copied}
}
