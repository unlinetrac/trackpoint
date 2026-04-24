package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/history"
	"github.com/user/trackpoint/internal/snapshot"
)

func makeStore(t *testing.T) *snapshot.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return snapshot.NewStore(dir)
}

func saveSnap(t *testing.T, store *snapshot.Store, label string, state map[string]string) *snapshot.Snapshot {
	t.Helper()
	snap := snapshot.New(label, state)
	if err := store.Save(snap); err != nil {
		t.Fatalf("save: %v", err)
	}
	return snap
}

func TestNewTimeline_LoadsSnapshots(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, "a", map[string]string{"k": "1"})
	s2 := saveSnap(t, store, "b", map[string]string{"k": "2"})

	tl, err := history.NewTimeline(store, []string{s1.ID, s2.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", tl.Len())
	}
}

func TestNewTimeline_DeduplicatesIDs(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, "a", map[string]string{"k": "1"})

	tl, err := history.NewTimeline(store, []string{s1.ID, s1.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl.Len() != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", tl.Len())
	}
}

func TestNewTimeline_MissingID_ReturnsError(t *testing.T) {
	store := makeStore(t)
	_, err := history.NewTimeline(store, []string{"nonexistent-id"})
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestTimeline_Pairs(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, "a", map[string]string{"k": "1"})
	s2 := saveSnap(t, store, "b", map[string]string{"k": "2"})
	s3 := saveSnap(t, store, "c", map[string]string{"k": "3"})

	tl, _ := history.NewTimeline(store, []string{s1.ID, s2.ID, s3.ID})
	pairs := tl.Pairs()
	if len(pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestTimeline_SortedByTime(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, "early", map[string]string{})
	time.Sleep(2 * time.Millisecond)
	s2 := saveSnap(t, store, "late", map[string]string{})

	// Load in reverse order
	tl, _ := history.NewTimeline(store, []string{s2.ID, s1.ID})
	sorted := tl.SortedByTime()

	if sorted.Entries()[0].Snapshot.Label != "early" {
		t.Errorf("expected first entry to be 'early', got %q", sorted.Entries()[0].Snapshot.Label)
	}
}
