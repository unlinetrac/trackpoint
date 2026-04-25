package annotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/trackpoint/internal/annotate"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "annotate-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	store := annotate.NewStore(tempDir(t))

	if err := store.Set("snap-001", "initial deploy"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	a, err := store.Get("snap-001")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if a.Note != "initial deploy" {
		t.Errorf("expected note %q, got %q", "initial deploy", a.Note)
	}
	if a.SnapshotID != "snap-001" {
		t.Errorf("expected snapshot_id %q, got %q", "snap-001", a.SnapshotID)
	}
}

func TestGet_NotFound_ReturnsEmpty(t *testing.T) {
	store := annotate.NewStore(tempDir(t))

	a, err := store.Get("nonexistent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if a.Note != "" {
		t.Errorf("expected empty note, got %q", a.Note)
	}
}

func TestDelete_RemovesAnnotation(t *testing.T) {
	dir := tempDir(t)
	store := annotate.NewStore(dir)

	_ = store.Set("snap-002", "to be removed")
	if err := store.Delete("snap-002"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	path := filepath.Join(dir, "snap-002.annotation.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted")
	}
}

func TestDelete_NotFound_NoError(t *testing.T) {
	store := annotate.NewStore(tempDir(t))
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("expected no error deleting nonexistent annotation, got %v", err)
	}
}

func TestSet_EmptyNote_DeletesAnnotation(t *testing.T) {
	dir := tempDir(t)
	store := annotate.NewStore(dir)

	_ = store.Set("snap-003", "will be cleared")
	if err := store.Set("snap-003", ""); err != nil {
		t.Fatalf("Set empty: %v", err)
	}

	a, err := store.Get("snap-003")
	if err != nil {
		t.Fatalf("Get after clear: %v", err)
	}
	if a.Note != "" {
		t.Errorf("expected empty note after clear, got %q", a.Note)
	}
}

func TestSet_EmptyID_ReturnsError(t *testing.T) {
	store := annotate.NewStore(tempDir(t))
	if err := store.Set("", "some note"); err == nil {
		t.Error("expected error for empty snapshot ID")
	}
}
