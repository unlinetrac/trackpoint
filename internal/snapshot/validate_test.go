package snapshot

import (
	"strings"
	"testing"
	"time"
)

func makeValidSnap() *Snapshot {
	return &Snapshot{
		ID:        "abc123",
		CreatedAt: time.Now(),
		Data:      map[string]string{"env": "prod"},
		Labels:    map[string]string{"team": "infra"},
	}
}

func TestValidate_ValidSnapshot_NoError(t *testing.T) {
	if err := Validate(makeValidSnap()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_NilSnapshot_ReturnsError(t *testing.T) {
	if err := Validate(nil); err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestValidate_EmptyID_ReturnsError(t *testing.T) {
	s := makeValidSnap()
	s.ID = "   "
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for blank id")
	}
	if !strings.Contains(err.Error(), "id must not be empty") {
		t.Errorf("unexpected message: %v", err)
	}
}

func TestValidate_ZeroCreatedAt_ReturnsError(t *testing.T) {
	s := makeValidSnap()
	s.CreatedAt = time.Time{}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for zero created_at")
	}
}

func TestValidate_NilData_ReturnsError(t *testing.T) {
	s := makeValidSnap()
	s.Data = nil
	if err := Validate(s); err == nil {
		t.Fatal("expected error for nil data")
	}
}

func TestValidate_OversizedValue_ReturnsError(t *testing.T) {
	s := makeValidSnap()
	s.Data["big"] = strings.Repeat("x", 4097)
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for oversized value")
	}
	if !strings.Contains(err.Error(), "exceeds 4096 bytes") {
		t.Errorf("unexpected message: %v", err)
	}
}

func TestValidate_BlankLabelValue_ReturnsError(t *testing.T) {
	s := makeValidSnap()
	s.Labels["env"] = "   "
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for blank label value")
	}
}

func TestValidate_MultipleProblems_ListsAll(t *testing.T) {
	s := makeValidSnap()
	s.ID = ""
	s.CreatedAt = time.Time{}
	ve, ok := Validate(s).(*ValidationError)
	if !ok {
		t.Fatal("expected *ValidationError")
	}
	if len(ve.Fields) < 2 {
		t.Errorf("expected at least 2 problems, got %d", len(ve.Fields))
	}
}
