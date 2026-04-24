package diff_test

import (
	"testing"

	"github.com/trackpoint/internal/diff"
	"github.com/trackpoint/internal/snapshot"
)

func makeSnapshot(state map[string]interface{}) *snapshot.Snapshot {
	return snapshot.New("test", state)
}

func TestCompare_NoChanges(t *testing.T) {
	state := map[string]interface{}{"replicas": 3, "image": "nginx:latest"}
	from := makeSnapshot(state)
	to := makeSnapshot(state)

	result := diff.Compare(from, to)
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestCompare_DetectsAddedKey(t *testing.T) {
	from := makeSnapshot(map[string]interface{}{"replicas": 2})
	to := makeSnapshot(map[string]interface{}{"replicas": 2, "image": "nginx:1.25"})

	result := diff.Compare(from, to)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added || result.Changes[0].Key != "image" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestCompare_DetectsRemovedKey(t *testing.T) {
	from := makeSnapshot(map[string]interface{}{"replicas": 2, "timeout": 30})
	to := makeSnapshot(map[string]interface{}{"replicas": 2})

	result := diff.Compare(from, to)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Removed || result.Changes[0].Key != "timeout" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestCompare_DetectsModifiedKey(t *testing.T) {
	from := makeSnapshot(map[string]interface{}{"replicas": 2})
	to := makeSnapshot(map[string]interface{}{"replicas": 5})

	result := diff.Compare(from, to)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != diff.Modified || c.Key != "replicas" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestResult_Summary(t *testing.T) {
	from := makeSnapshot(map[string]interface{}{"a": 1})
	to := makeSnapshot(map[string]interface{}{"a": 2})

	result := diff.Compare(from, to)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
