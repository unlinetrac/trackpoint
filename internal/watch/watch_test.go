package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/snapshot"
	"github.com/user/trackpoint/internal/watch"
)

func makeStore(t *testing.T) *snapshot.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "snapshots")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	return snapshot.NewStore(dir)
}

func saveSnap(t *testing.T, store *snapshot.Store, data map[string]string) *snapshot.Snapshot {
	t.Helper()
	snap := snapshot.New("test", data)
	if err := store.Save(snap); err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestWatch_EmitsChangeWhenDifferent(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, map[string]string{"key": "a"})
	s2 := saveSnap(t, store, map[string]string{"key": "b"})

	w := watch.New(store, watch.Config{Interval: time.Millisecond, MaxChecks: 1})
	out := make(chan watch.Change, 10)

	if err := w.Watch([]string{s1.ID, s2.ID}, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) == 0 {
		t.Fatal("expected at least one change event")
	}
	change := <-out
	if len(change.Result.Changes) == 0 {
		t.Error("expected changes in result")
	}
}

func TestWatch_NoEmitWhenIdentical(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, map[string]string{"key": "same"})
	s2 := saveSnap(t, store, map[string]string{"key": "same"})

	w := watch.New(store, watch.Config{Interval: time.Millisecond, MaxChecks: 1})
	out := make(chan watch.Change, 10)

	if err := w.Watch([]string{s1.ID, s2.ID}, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 0 {
		t.Errorf("expected no change events, got %d", len(out))
	}
}

func TestWatch_SingleSnapshot_NoChanges(t *testing.T) {
	store := makeStore(t)
	s1 := saveSnap(t, store, map[string]string{"key": "val"})

	w := watch.New(store, watch.Config{Interval: time.Millisecond, MaxChecks: 1})
	out := make(chan watch.Change, 10)

	if err := w.Watch([]string{s1.ID}, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 0 {
		t.Errorf("expected no change events for single snapshot, got %d", len(out))
	}
}
