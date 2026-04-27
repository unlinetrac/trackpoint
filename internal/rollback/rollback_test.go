package rollback_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/rollback"
	"github.com/user/trackpoint/internal/snapshot"
)

// fakeStore implements rollback.SnapshotLoader.
type fakeStore struct {
	snaps map[string]*snapshot.Snapshot
	order []string
}

func (f *fakeStore) Load(id string) (*snapshot.Snapshot, error) {
	s, ok := f.snaps[id]
	if !ok {
		return nil, fmt.Errorf("not found: %s", id)
	}
	return s, nil
}

func (f *fakeStore) List() ([]string, error) {
	return f.order, nil
}

func makeSnap(id string, data map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Data:      data,
	}
}

func TestFind_ReturnsPreviousSnapshot(t *testing.T) {
	store := &fakeStore{
		order: []string{"snap-1", "snap-2", "snap-3"},
		snaps: map[string]*snapshot.Snapshot{
			"snap-1": makeSnap("snap-1", map[string]string{"a": "1"}),
			"snap-2": makeSnap("snap-2", map[string]string{"a": "1", "b": "2"}),
			"snap-3": makeSnap("snap-3", map[string]string{"a": "1", "b": "2", "c": "3"}),
		},
	}
	res, err := rollback.Find(store, "snap-3", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Target.ID != "snap-2" {
		t.Errorf("expected snap-2, got %s", res.Target.ID)
	}
}

func TestFind_NoDiff_SkipsChangedSnapshots(t *testing.T) {
	store := &fakeStore{
		order: []string{"snap-1", "snap-2", "snap-3"},
		snaps: map[string]*snapshot.Snapshot{
			"snap-1": makeSnap("snap-1", map[string]string{"a": "1"}),
			"snap-2": makeSnap("snap-2", map[string]string{"a": "1", "b": "2"}),
			"snap-3": makeSnap("snap-3", map[string]string{"a": "1"}),
		},
	}
	res, err := rollback.Find(store, "snap-3", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Target.ID != "snap-1" {
		t.Errorf("expected snap-1, got %s", res.Target.ID)
	}
}

func TestFind_NoCandidate_ReturnsError(t *testing.T) {
	store := &fakeStore{
		order: []string{"snap-1"},
		snaps: map[string]*snapshot.Snapshot{
			"snap-1": makeSnap("snap-1", map[string]string{"x": "y"}),
		},
	}
	_, err := rollback.Find(store, "snap-1", false)
	if err != rollback.ErrNoCandidate {
		t.Errorf("expected ErrNoCandidate, got %v", err)
	}
}
