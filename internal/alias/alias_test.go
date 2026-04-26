package alias_test

import (
	"errors"
	"os"
	"testing"

	"github.com/yourorg/trackpoint/internal/alias"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "alias-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	if err := st.Set("prod", "abc123"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	id, err := st.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if id != "abc123" {
		t.Errorf("got %q, want %q", id, "abc123")
	}
}

func TestGet_NotFound(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	_, err := st.Get("missing")
	if !errors.Is(err, alias.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSet_InvalidAlias(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	if err := st.Set("bad name!", "abc"); err == nil {
		t.Error("expected error for invalid alias name")
	}
}

func TestSet_EmptySnapshotID(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	if err := st.Set("ok", ""); err == nil {
		t.Error("expected error for empty snapshot ID")
	}
}

func TestDelete_RemovesAlias(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	_ = st.Set("staging", "xyz789")
	if err := st.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := st.Get("staging")
	if !errors.Is(err, alias.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDelete_NotFound_NoError(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	if err := st.Delete("nonexistent"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestList_ReturnsAllAliases(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	_ = st.Set("alpha", "id1")
	_ = st.Set("beta", "id2")
	names, err := st.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(names))
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	names, err := st.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}
