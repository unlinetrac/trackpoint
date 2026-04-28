package merge_test

import (
	"testing"
	"time"

	"github.com/your-org/trackpoint/internal/merge"
	"github.com/your-org/trackpoint/internal/snapshot"
)

func makeSnap(id string, data map[string]string) *snapshot.Snapshot {
	snap := &snapshot.Snapshot{
		ID:        id,
		Timestamp: time.Now(),
		Data:      data,
	}
	return snap
}

func TestRun_MergesDisjointKeys(t *testing.T) {
	a := makeSnap("a", map[string]string{
		"host": "web-01",
		"env":  "prod",
	})
	b := makeSnap("b", map[string]string{
		"region": "us-east-1",
		"tier":   "frontend",
	})

	result, err := merge.Run(a, b, merge.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Data["host"] != "web-01" {
		t.Errorf("expected host=web-01, got %q", result.Data["host"])
	}
	if result.Data["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", result.Data["region"])
	}
	if len(result.Data) != 4 {
		t.Errorf("expected 4 keys, got %d", len(result.Data))
	}
}

func TestRun_ConflictResolution_LastWins(t *testing.T) {
	a := makeSnap("a", map[string]string{
		"host": "web-01",
		"env":  "staging",
	})
	b := makeSnap("b", map[string]string{
		"host": "web-02",
		"env":  "prod",
	})

	result, err := merge.Run(a, b, merge.Options{Strategy: merge.StrategyLastWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Data["host"] != "web-02" {
		t.Errorf("expected host=web-02 (last wins), got %q", result.Data["host"])
	}
	if result.Data["env"] != "prod" {
		t.Errorf("expected env=prod (last wins), got %q", result.Data["env"])
	}
}

func TestRun_ConflictResolution_FirstWins(t *testing.T) {
	a := makeSnap("a", map[string]string{
		"host": "web-01",
	})
	b := makeSnap("b", map[string]string{
		"host": "web-02",
	})

	result, err := merge.Run(a, b, merge.Options{Strategy: merge.StrategyFirstWins})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Data["host"] != "web-01" {
		t.Errorf("expected host=web-01 (first wins), got %q", result.Data["host"])
	}
}

func TestRun_ConflictResolution_ErrorOnConflict(t *testing.T) {
	a := makeSnap("a", map[string]string{
		"host": "web-01",
	})
	b := makeSnap("b", map[string]string{
		"host": "web-02",
	})

	_, err := merge.Run(a, b, merge.Options{Strategy: merge.StrategyError})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestRun_NilSnapshotA_ReturnsError(t *testing.T) {
	b := makeSnap("b", map[string]string{"key": "val"})

	_, err := merge.Run(nil, b, merge.Options{})
	if err == nil {
		t.Fatal("expected error for nil snapshot a")
	}
}

func TestRun_NilSnapshotB_ReturnsError(t *testing.T) {
	a := makeSnap("a", map[string]string{"key": "val"})

	_, err := merge.Run(a, nil, merge.Options{})
	if err == nil {
		t.Fatal("expected error for nil snapshot b")
	}
}

func TestRun_EmptySnapshots_ReturnsEmptyData(t *testing.T) {
	a := makeSnap("a", map[string]string{})
	b := makeSnap("b", map[string]string{})

	result, err := merge.Run(a, b, merge.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 0 {
		t.Errorf("expected empty data, got %d keys", len(result.Data))
	}
}
