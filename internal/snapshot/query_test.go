package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/trackpoint/internal/snapshot"
)

type fakeStore struct {
	snaps []*snapshot.Snapshot
}

func (f *fakeStore) List() ([]string, error) {
	ids := make([]string, len(f.snaps))
	for i, s := range f.snaps {
		ids[i] = s.ID
	}
	return ids, nil
}

func (f *fakeStore) Load(id string) (*snapshot.Snapshot, error) {
	for _, s := range f.snaps {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, snapshot.ErrNotFound
}

func makeQuerySnap(id, label string, tags map[string]string, at time.Time) *snapshot.Snapshot {
	s, _ := snapshot.New(label, map[string]string{"k": "v"}, tags)
	s.ID = id
	s.CreatedAt = at
	return s
}

func TestQuery_NoOptions_ReturnsAll(t *testing.T) {
	now := time.Now()
	store := &fakeStore{snaps: []*snapshot.Snapshot{
		makeQuerySnap("a", "alpha", nil, now),
		makeQuerySnap("b", "beta", nil, now),
	}}
	results, err := snapshot.Query(store, snapshot.QueryOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestQuery_LabelFilter(t *testing.T) {
	now := time.Now()
	store := &fakeStore{snaps: []*snapshot.Snapshot{
		makeQuerySnap("a", "production-deploy", nil, now),
		makeQuerySnap("b", "staging-deploy", nil, now),
		makeQuerySnap("c", "rollback", nil, now),
	}}
	results, err := snapshot.Query(store, snapshot.QueryOptions{Label: "deploy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestQuery_SinceFilter(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	cutoff := base.Add(24 * time.Hour)
	store := &fakeStore{snaps: []*snapshot.Snapshot{
		makeQuerySnap("old", "old", nil, base),
		makeQuerySnap("new", "new", nil, base.Add(48*time.Hour)),
	}}
	results, err := snapshot.Query(store, snapshot.QueryOptions{Since: &cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].ID != "new" {
		t.Fatalf("expected only 'new', got %v", results)
	}
}

func TestQuery_TagFilter(t *testing.T) {
	now := time.Now()
	store := &fakeStore{snaps: []*snapshot.Snapshot{
		makeQuerySnap("a", "a", map[string]string{"env": "prod"}, now),
		makeQuerySnap("b", "b", map[string]string{"env": "staging"}, now),
	}}
	results, err := snapshot.Query(store, snapshot.QueryOptions{Tags: map[string]string{"env": "prod"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].ID != "a" {
		t.Fatalf("expected only snap 'a', got %v", results)
	}
}

func TestQuery_Limit(t *testing.T) {
	now := time.Now()
	store := &fakeStore{snaps: []*snapshot.Snapshot{
		makeQuerySnap("a", "a", nil, now),
		makeQuerySnap("b", "b", nil, now),
		makeQuerySnap("c", "c", nil, now),
	}}
	results, err := snapshot.Query(store, snapshot.QueryOptions{Limit: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}
