package baseline_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/baseline"
	"github.com/user/trackpoint/internal/snapshot"
)

// fakeStore is an in-memory snapshot store for testing.
type fakeStore struct {
	data map[string]*snapshot.Snapshot
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]*snapshot.Snapshot)}
}

func (f *fakeStore) Save(s *snapshot.Snapshot) error {
	f.data[s.ID] = s
	return nil
}

func (f *fakeStore) Load(id string) (*snapshot.Snapshot, error) {
	s, ok := f.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func makeSnap(data map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "snap-abc",
		CreatedAt: time.Now(),
		Data:      data,
	}
}

func TestSet_And_Get_RoundTrip(t *testing.T) {
	bs := baseline.New(newFakeStore())
	s := makeSnap(map[string]string{"env": "prod"})

	if err := bs.Set(s); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := bs.Get()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Data["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", got.Data["env"])
	}
}

func TestGet_NoBaseline_ReturnsErrNoBaseline(t *testing.T) {
	bs := baseline.New(newFakeStore())
	_, err := bs.Get()
	if !errors.Is(err, baseline.ErrNoBaseline) {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestSet_NilSnapshot_ReturnsError(t *testing.T) {
	bs := baseline.New(newFakeStore())
	if err := bs.Set(nil); err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestIsSet_ReturnsTrueAfterSet(t *testing.T) {
	bs := baseline.New(newFakeStore())
	if bs.IsSet() {
		t.Error("expected IsSet=false before any Set")
	}
	_ = bs.Set(makeSnap(map[string]string{"k": "v"}))
	if !bs.IsSet() {
		t.Error("expected IsSet=true after Set")
	}
}

func TestClear_RemovesBaseline(t *testing.T) {
	bs := baseline.New(newFakeStore())
	_ = bs.Set(makeSnap(map[string]string{"k": "v"}))
	_ = bs.Clear()
	if bs.IsSet() {
		t.Error("expected IsSet=false after Clear")
	}
}
