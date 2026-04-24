package snapshot_test

import (
	"os"
	"testing"

	"github.com/user/trackpoint/internal/snapshot"
)

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	snap := snapshot.New("env:prod", map[string]string{
		"DB_HOST": "localhost",
		"PORT":    "5432",
	})

	if err := store.Save(snap); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := store.Load(snap.ID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.ID != snap.ID {
		t.Errorf("ID mismatch: got %q, want %q", loaded.ID, snap.ID)
	}
	if loaded.Source != snap.Source {
		t.Errorf("Source mismatch: got %q, want %q", loaded.Source, snap.Source)
	}
	if loaded.State["DB_HOST"] != snap.State["DB_HOST"] {
		t.Errorf("State mismatch for DB_HOST")
	}
}

func TestStore_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestStore_List_Empty(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	ids, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected empty list, got %v", ids)
	}
}

func TestStore_List_ReturnsSortedIDs(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	for _, source := range []string{"env:staging", "env:prod", "env:dev"} {
		snap := snapshot.New(source, map[string]string{"K": "V"})
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	ids, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 IDs, got %d", len(ids))
	}
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Errorf("IDs not sorted: %v", ids)
		}
	}
}

func TestStore_List_NoDir(t *testing.T) {
	store := snapshot.NewStore("/tmp/trackpoint_nonexistent_dir_xyz")
	ids, err := store.List()
	if err != nil {
		t.Fatalf("List() on missing dir should not error, got %v", err)
	}
	if ids != nil {
		t.Errorf("expected nil slice, got %v", ids)
	}
	_ = os.RemoveAll("/tmp/trackpoint_nonexistent_dir_xyz")
}
