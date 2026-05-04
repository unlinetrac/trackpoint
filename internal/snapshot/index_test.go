package snapshot

import (
	"testing"
	"time"
)

func makeIndexSnap(id string, labels map[string]string) *Snapshot {
	return &Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Labels:    labels,
		Data:      map[string]string{"key": "value"},
	}
}

func TestNewIndex_Size(t *testing.T) {
	snaps := []*Snapshot{
		makeIndexSnap("a", map[string]string{"env": "prod"}),
		makeIndexSnap("b", map[string]string{"env": "staging"}),
		makeIndexSnap("c", map[string]string{"env": "prod", "region": "us-east"}),
	}
	idx := NewIndex(snaps)
	if idx.Size() != 3 {
		t.Fatalf("expected size 3, got %d", idx.Size())
	}
}

func TestNewIndex_NilSnapshotsIgnored(t *testing.T) {
	snaps := []*Snapshot{nil, makeIndexSnap("a", nil), nil}
	idx := NewIndex(snaps)
	if idx.Size() != 1 {
		t.Fatalf("expected size 1, got %d", idx.Size())
	}
}

func TestFindByLabel_ReturnsMatchingSnapshots(t *testing.T) {
	snaps := []*Snapshot{
		makeIndexSnap("a", map[string]string{"env": "prod"}),
		makeIndexSnap("b", map[string]string{"env": "staging"}),
		makeIndexSnap("c", map[string]string{"env": "prod"}),
	}
	idx := NewIndex(snaps)
	results := idx.FindByLabel("env", "prod")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFindByLabel_NoMatch_ReturnsEmpty(t *testing.T) {
	idx := NewIndex([]*Snapshot{
		makeIndexSnap("a", map[string]string{"env": "prod"}),
	})
	results := idx.FindByLabel("env", "dev")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestFindByID_ReturnsSnapshot(t *testing.T) {
	s := makeIndexSnap("abc123", nil)
	idx := NewIndex([]*Snapshot{s})
	found := idx.FindByID("abc123")
	if found == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if found.ID != "abc123" {
		t.Fatalf("expected ID abc123, got %s", found.ID)
	}
}

func TestFindByID_NotFound_ReturnsNil(t *testing.T) {
	idx := NewIndex([]*Snapshot{})
	if idx.FindByID("missing") != nil {
		t.Fatal("expected nil for unknown ID")
	}
}

func TestFindByLabelKey_ReturnsAllValues(t *testing.T) {
	snaps := []*Snapshot{
		makeIndexSnap("a", map[string]string{"env": "prod"}),
		makeIndexSnap("b", map[string]string{"env": "staging"}),
		makeIndexSnap("c", map[string]string{"region": "us-east"}),
	}
	idx := NewIndex(snaps)
	results := idx.FindByLabelKey("env")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFindByLabelKey_NoMatch_ReturnsEmpty(t *testing.T) {
	idx := NewIndex([]*Snapshot{
		makeIndexSnap("a", map[string]string{"env": "prod"}),
	})
	results := idx.FindByLabelKey("nonexistent")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
