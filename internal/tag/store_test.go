package tag

import (
	"os"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "tagstore-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestTagStore_SaveAndLoad(t *testing.T) {
	store := NewTagStore(tempDir(t))
	tags := Tags{"env": "prod", "team": "platform"}

	if err := store.Save("snap1", tags); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load("snap1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", loaded["env"])
	}
	if loaded["team"] != "platform" {
		t.Errorf("expected team=platform, got %s", loaded["team"])
	}
}

func TestTagStore_Load_NotFound_ReturnsEmpty(t *testing.T) {
	store := NewTagStore(tempDir(t))
	tags, err := store.Load("nonexistent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestTagStore_Delete(t *testing.T) {
	store := NewTagStore(tempDir(t))
	if err := store.Save("snap2", Tags{"x": "y"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if err := store.Delete("snap2"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	tags, err := store.Load("snap2")
	if err != nil {
		t.Fatalf("Load after delete failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty after delete, got %v", tags)
	}
}

func TestTagStore_Delete_NotFound_NoError(t *testing.T) {
	store := NewTagStore(tempDir(t))
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("expected no error deleting nonexistent tag file, got %v", err)
	}
}
