package validate

import (
	"testing"
	"time"

	"github.com/yourorg/trackpoint/internal/snapshot"
)

func makeSnap(id string, data map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Data:      data,
		Labels:    map[string]string{},
	}
}

func TestRun_NilSnapshot_ReturnsError(t *testing.T) {
	_, err := Run(nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRun_ValidCleanSnapshot_NoIssues(t *testing.T) {
	snap := makeSnap("abc", map[string]string{"region": "us-east-1"})
	res, err := Run(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(res.Issues), res.Issues)
	}
	if res.HasErrors() {
		t.Error("HasErrors should be false for a clean snapshot")
	}
}

func TestRun_EmptyID_ReturnsErrorIssue(t *testing.T) {
	snap := makeSnap("", map[string]string{"k": "v"})
	res, err := Run(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasErrors() {
		t.Error("expected at least one error issue")
	}
}

func TestRun_EmptyValue_ReturnsWarning(t *testing.T) {
	snap := makeSnap("id1", map[string]string{"key": ""})
	res, err := Run(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var warnings int
	for _, i := range res.Issues {
		if i.Level == "warning" {
			warnings++
		}
	}
	if warnings == 0 {
		t.Error("expected at least one warning for empty value")
	}
}

func TestRun_SnapshotIDPropagated(t *testing.T) {
	snap := makeSnap("snap-42", map[string]string{"x": "y"})
	res, _ := Run(snap)
	if res.SnapshotID != "snap-42" {
		t.Errorf("expected SnapshotID snap-42, got %s", res.SnapshotID)
	}
}
