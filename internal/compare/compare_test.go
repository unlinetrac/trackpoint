package compare_test

import (
	"testing"
	"time"

	"github.com/user/trackpoint/internal/compare"
	"github.com/user/trackpoint/internal/snapshot"
)

func makeSnap(id string, data map[string]string) *snapshot.Snapshot {
	s := &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Data:      data,
	}
	return s
}

func TestRun_TooFewSnapshots(t *testing.T) {
	_, err := compare.Run([]*snapshot.Snapshot{makeSnap("a", map[string]string{"k": "v"})})
	if err == nil {
		t.Fatal("expected error for single snapshot, got nil")
	}
}

func TestRun_NoChanges(t *testing.T) {
	data := map[string]string{"env": "prod", "region": "us-east-1"}
	snaps := []*snapshot.Snapshot{
		makeSnap("snap1", data),
		makeSnap("snap2", data),
	}
	res, err := compare.Run(snaps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasChanges() {
		t.Error("expected no changes, but HasChanges returned true")
	}
	if res.SnapshotCount != 2 {
		t.Errorf("expected SnapshotCount=2, got %d", res.SnapshotCount)
	}
}

func TestRun_DetectsNetChanges(t *testing.T) {
	from := makeSnap("snap1", map[string]string{"a": "1", "b": "2"})
	mid := makeSnap("snap2", map[string]string{"a": "1", "b": "99", "c": "3"})
	to := makeSnap("snap3", map[string]string{"a": "1", "c": "3"})

	res, err := compare.Run([]*snapshot.Snapshot{from, mid, to})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasChanges() {
		t.Fatal("expected changes, got none")
	}
	if len(res.Diff.Removed) != 1 {
		t.Errorf("expected 1 removed key, got %d", len(res.Diff.Removed))
	}
	if res.FromID != "snap1" || res.ToID != "snap3" {
		t.Errorf("unexpected from/to IDs: %s / %s", res.FromID, res.ToID)
	}
	if res.SnapshotCount != 3 {
		t.Errorf("expected SnapshotCount=3, got %d", res.SnapshotCount)
	}
}

func TestRun_UsesFirstAndLast(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap("first", map[string]string{"x": "old"}),
		makeSnap("middle", map[string]string{"x": "mid"}),
		makeSnap("last", map[string]string{"x": "new"}),
	}
	res, err := compare.Run(snaps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FromID != "first" {
		t.Errorf("expected FromID=first, got %s", res.FromID)
	}
	if res.ToID != "last" {
		t.Errorf("expected ToID=last, got %s", res.ToID)
	}
	if len(res.Diff.Modified) != 1 {
		t.Errorf("expected 1 modified key, got %d", len(res.Diff.Modified))
	}
}
