package snapshot

import (
	"testing"
)

func makePatchSnap(data map[string]string) *Snapshot {
	s, _ := New(data, nil)
	return s
}

func TestPatch_SetAddsKey(t *testing.T) {
	s := makePatchSnap(map[string]string{"env": "prod"})
	ops := []PatchOp{{Key: "region", Op: "set", Value: "us-east-1"}}

	out, err := Patch(s, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", out.Data["region"])
	}
	// Original must not be mutated.
	if _, ok := s.Data["region"]; ok {
		t.Error("original snapshot was mutated")
	}
}

func TestPatch_DeleteRemovesKey(t *testing.T) {
	s := makePatchSnap(map[string]string{"env": "prod", "debug": "true"})
	ops := []PatchOp{{Key: "debug", Op: "delete"}}

	out, err := Patch(s, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Data["debug"]; ok {
		t.Error("expected debug key to be deleted")
	}
	if out.Data["env"] != "prod" {
		t.Errorf("expected env=prod to remain, got %q", out.Data["env"])
	}
}

func TestPatch_RenameMovesKey(t *testing.T) {
	s := makePatchSnap(map[string]string{"old_name": "value"})
	ops := []PatchOp{{Key: "old_name", Op: "rename", NewKey: "new_name"}}

	out, err := Patch(s, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Data["old_name"]; ok {
		t.Error("expected old_name to be removed")
	}
	if out.Data["new_name"] != "value" {
		t.Errorf("expected new_name=value, got %q", out.Data["new_name"])
	}
}

func TestPatch_IDChangesAfterPatch(t *testing.T) {
	s := makePatchSnap(map[string]string{"k": "v"})
	ops := []PatchOp{{Key: "k", Op: "set", Value: "changed"}}

	out, err := Patch(s, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID == s.ID {
		t.Error("expected ID to change after patch")
	}
}

func TestPatch_NilSnapshot_ReturnsError(t *testing.T) {
	_, err := Patch(nil, []PatchOp{{Key: "k", Op: "set", Value: "v"}})
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestPatch_UnknownOp_ReturnsError(t *testing.T) {
	s := makePatchSnap(map[string]string{"k": "v"})
	_, err := Patch(s, []PatchOp{{Key: "k", Op: "upsert"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatch_RenameSourceMissing_ReturnsError(t *testing.T) {
	s := makePatchSnap(map[string]string{"k": "v"})
	_, err := Patch(s, []PatchOp{{Key: "missing", Op: "rename", NewKey: "other"}})
	if err == nil {
		t.Error("expected error when rename source key is missing")
	}
}
