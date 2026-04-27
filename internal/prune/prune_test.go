package prune_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/user/trackpoint/internal/prune"
)

// fakeStore is an in-memory Store implementation for testing.
type fakeStore struct {
	ids       []string
	deleted   []string
	listErr   error
	deleteErr error
}

func (f *fakeStore) List() ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	out := make([]string, len(f.ids))
	copy(out, f.ids)
	sort.Strings(out)
	return out, nil
}

func (f *fakeStore) Delete(id string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	f.deleted = append(f.deleted, id)
	return nil
}

func TestRun_KeepLast_DeletesOldest(t *testing.T) {
	store := &fakeStore{ids: []string{"a", "b", "c", "d", "e"}}
	result, err := prune.Run(store, prune.Options{KeepLast: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Deleted) != 3 {
		t.Errorf("expected 3 deleted, got %d", len(result.Deleted))
	}
	if len(result.Skipped) != 2 {
		t.Errorf("expected 2 skipped, got %d", len(result.Skipped))
	}
	if len(store.deleted) != 3 {
		t.Errorf("expected 3 actual deletes, got %d", len(store.deleted))
	}
}

func TestRun_FewerThanKeepLast_NothingDeleted(t *testing.T) {
	store := &fakeStore{ids: []string{"a", "b"}}
	result, err := prune.Run(store, prune.Options{KeepLast: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Deleted) != 0 {
		t.Errorf("expected 0 deleted, got %d", len(result.Deleted))
	}
	if len(store.deleted) != 0 {
		t.Errorf("expected 0 actual deletes, got %d", len(store.deleted))
	}
}

func TestRun_DryRun_DoesNotDelete(t *testing.T) {
	store := &fakeStore{ids: []string{"a", "b", "c"}}
	result, err := prune.Run(store, prune.Options{KeepLast: 1, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Deleted) != 2 {
		t.Errorf("expected 2 in deleted list, got %d", len(result.Deleted))
	}
	if len(store.deleted) != 0 {
		t.Errorf("dry run should not delete; got %d deletes", len(store.deleted))
	}
}

func TestRun_InvalidKeepLast_ReturnsError(t *testing.T) {
	store := &fakeStore{ids: []string{"a"}}
	_, err := prune.Run(store, prune.Options{KeepLast: 0})
	if err == nil {
		t.Fatal("expected error for KeepLast=0, got nil")
	}
}

func TestRun_ListError_Propagates(t *testing.T) {
	store := &fakeStore{listErr: errors.New("disk failure")}
	_, err := prune.Run(store, prune.Options{KeepLast: 1})
	if err == nil {
		t.Fatal("expected error from List, got nil")
	}
}
