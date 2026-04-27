package search_test

import (
	"testing"
	"time"

	"github.com/trackpoint/internal/search"
	"github.com/trackpoint/internal/snapshot"
)

func makeSnap(id string, state map[string]string, tags map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		State:     state,
		Tags:      tags,
	}
}

func TestRun_NoOptions_ReturnsAll(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap("a", map[string]string{"key": "val"}, nil),
		makeSnap("b", map[string]string{"other": "thing"}, nil),
	}
	results := search.Run(snaps, search.Options{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRun_KeyContains_FiltersCorrectly(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap("a", map[string]string{"db.host": "localhost", "db.port": "5432"}, nil),
		makeSnap("b", map[string]string{"app.name": "trackpoint"}, nil),
	}
	results := search.Run(snaps, search.Options{KeyContains: "db."})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Snapshot.ID != "a" {
		t.Errorf("expected snapshot a, got %s", results[0].Snapshot.ID)
	}
	if len(results[0].MatchedKeys) != 2 {
		t.Errorf("expected 2 matched keys, got %d", len(results[0].MatchedKeys))
	}
}

func TestRun_ValueContains_FiltersCorrectly(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap("a", map[string]string{"host": "prod-server"}, nil),
		makeSnap("b", map[string]string{"host": "staging-server"}, nil),
	}
	results := search.Run(snaps, search.Options{ValueContains: "prod"})
	if len(results) != 1 || results[0].Snapshot.ID != "a" {
		t.Errorf("expected only snapshot a")
	}
}

func TestRun_TagFilter_FiltersCorrectly(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap("a", map[string]string{"k": "v"}, map[string]string{"env": "prod"}),
		makeSnap("b", map[string]string{"k": "v"}, map[string]string{"env": "staging"}),
	}
	results := search.Run(snaps, search.Options{Tags: map[string]string{"env": "prod"}})
	if len(results) != 1 || results[0].Snapshot.ID != "a" {
		t.Errorf("expected only snapshot a")
	}
}

func TestRun_NilSnapshot_Skipped(t *testing.T) {
	snaps := []*snapshot.Snapshot{nil, makeSnap("a", map[string]string{"k": "v"}, nil)}
	results := search.Run(snaps, search.Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}
