package replay_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/replay"
	"github.com/user/trackpoint/internal/snapshot"
)

func makeSnap(t *testing.T, data map[string]string) *snapshot.Snapshot {
	t.Helper()
	s, err := snapshot.New("test", data, nil)
	if err != nil {
		t.Fatalf("makeSnap: %v", err)
	}
	s.CreatedAt = time.Now()
	return s
}

func TestRun_TooFewSnapshots(t *testing.T) {
	err := replay.Run([]*snapshot.Snapshot{makeSnap(t, nil)}, replay.Options{}, func(replay.Frame) error { return nil })
	if err == nil {
		t.Fatal("expected error for fewer than two snapshots")
	}
}

func TestRun_EmitsFramesInOrder(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap(t, map[string]string{"a": "1"}),
		makeSnap(t, map[string]string{"a": "2"}),
		makeSnap(t, map[string]string{"a": "3"}),
	}

	var indices []int
	err := replay.Run(snaps, replay.Options{}, func(f replay.Frame) error {
		indices = append(indices, f.Index)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(indices) != 2 || indices[0] != 1 || indices[1] != 2 {
		t.Fatalf("unexpected frame indices: %v", indices)
	}
}

func TestRun_OnlyChanged_SkipsIdenticalFrames(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap(t, map[string]string{"a": "1"}),
		makeSnap(t, map[string]string{"a": "1"}),
		makeSnap(t, map[string]string{"a": "2"}),
	}

	var count int
	err := replay.Run(snaps, replay.Options{OnlyChanged: true}, func(f replay.Frame) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 frame, got %d", count)
	}
}

func TestRun_CallbackError_HaltsReplay(t *testing.T) {
	snaps := []*snapshot.Snapshot{
		makeSnap(t, map[string]string{"a": "1"}),
		makeSnap(t, map[string]string{"a": "2"}),
		makeSnap(t, map[string]string{"a": "3"}),
	}

	sentinel := errors.New("stop")
	var count int
	err := replay.Run(snaps, replay.Options{}, func(f replay.Frame) error {
		count++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected callback to be called once, got %d", count)
	}
}
