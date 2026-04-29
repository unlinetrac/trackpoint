package pipeline_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/trackpoint/internal/pipeline"
	"github.com/user/trackpoint/internal/snapshot"
)

func makeSnap(label string, data map[string]string) *snapshot.Snapshot {
	return snapshot.New(label, data)
}

func TestRun_EmptyPipeline_ReturnsInputUnchanged(t *testing.T) {
	p := pipeline.New()
	snap := makeSnap("s1", map[string]string{"k": "v"})
	out, err := p.Run(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data["k"] != "v" {
		t.Errorf("expected data unchanged, got %v", out.Data)
	}
}

func TestRun_NilInput_ReturnsError(t *testing.T) {
	p := pipeline.New()
	_, err := p.Run(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestRun_StageError_HaltsPipeline(t *testing.T) {
	sentinel := errors.New("boom")
	p := pipeline.New().
		Add("fail", func(s *snapshot.Snapshot) (*snapshot.Snapshot, error) {
			return nil, sentinel
		})
	_, err := p.Run(makeSnap("s", nil))
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected stage name in error, got %v", err)
	}
}

func TestRun_MultipleStages_AppliedInOrder(t *testing.T) {
	p := pipeline.New().
		Add("add-a", func(s *snapshot.Snapshot) (*snapshot.Snapshot, error) {
			d := map[string]string{"a": "1"}
			for k, v := range s.Data {
				d[k] = v
			}
			return snapshot.New(s.Label, d), nil
		}).
		Add("add-b", func(s *snapshot.Snapshot) (*snapshot.Snapshot, error) {
			d := map[string]string{"b": "2"}
			for k, v := range s.Data {
				d[k] = v
			}
			return snapshot.New(s.Label, d), nil
		})
	out, err := p.Run(makeSnap("s", map[string]string{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data["a"] != "1" || out.Data["b"] != "2" {
		t.Errorf("unexpected data: %v", out.Data)
	}
}

func TestStages_ReturnsCopy(t *testing.T) {
	p := pipeline.New().Add("noop", func(s *snapshot.Snapshot) (*snapshot.Snapshot, error) {
		return s, nil
	})
	stages := p.Stages()
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
}
