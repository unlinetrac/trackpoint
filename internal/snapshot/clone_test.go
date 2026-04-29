package snapshot

import (
	"testing"
)

func TestClone_ReturnsDeepCopy(t *testing.T) {
	orig := New(map[string]string{"host": "web-01", "env": "prod"})
	orig.Tags = map[string]string{"team": "infra"}

	cloned, err := Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cloned.ID == orig.ID {
		t.Errorf("expected cloned ID to differ from original, got same: %s", cloned.ID)
	}
	if cloned.Data["host"] != orig.Data["host"] {
		t.Errorf("expected cloned data to match original")
	}
	if cloned.Tags["team"] != orig.Tags["team"] {
		t.Errorf("expected cloned tags to match original")
	}

	// Mutating clone should not affect original
	cloned.Data["host"] = "web-99"
	if orig.Data["host"] != "web-01" {
		t.Errorf("mutating clone affected original")
	}
}

func TestClone_NilInput_ReturnsError(t *testing.T) {
	_, err := Clone(nil)
	if err == nil {
		t.Fatal("expected error for nil input, got nil")
	}
}

func TestCloneWithOverrides_MergesKeys(t *testing.T) {
	orig := New(map[string]string{"host": "web-01", "env": "prod"})

	cloned, err := CloneWithOverrides(orig, map[string]string{"env": "staging", "region": "us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cloned.Data["env"] != "staging" {
		t.Errorf("expected env=staging, got %s", cloned.Data["env"])
	}
	if cloned.Data["host"] != "web-01" {
		t.Errorf("expected host=web-01 to be preserved, got %s", cloned.Data["host"])
	}
	if cloned.Data["region"] != "us-east" {
		t.Errorf("expected region=us-east to be added, got %s", cloned.Data["region"])
	}
	if orig.Data["env"] != "prod" {
		t.Errorf("override mutated original")
	}
}

func TestCloneWithOverrides_IDChangesWithNewData(t *testing.T) {
	orig := New(map[string]string{"a": "1"})

	cloned, err := CloneWithOverrides(orig, map[string]string{"a": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cloned.ID == orig.ID {
		t.Errorf("expected ID to change after override, but got same ID")
	}
}

func TestCloneWithOverrides_NilInput_ReturnsError(t *testing.T) {
	_, err := CloneWithOverrides(nil, map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error for nil snapshot, got nil")
	}
}
