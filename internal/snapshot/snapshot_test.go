package snapshot_test

import (
	"testing"

	"github.com/yourorg/trackpoint/internal/snapshot"
)

func TestNew_SetsFieldsCorrectly(t *testing.T) {
	entries := map[string]string{
		"service/api/replicas": "3",
		"service/api/image":    "nginx:1.25",
	}
	s := snapshot.New("pre-deploy", entries)

	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if s.Label != "pre-deploy" {
		t.Errorf("expected label 'pre-deploy', got %q", s.Label)
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.Entries))
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	original := snapshot.New("test", map[string]string{
		"db/host": "localhost",
		"db/port": "5432",
	})

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	restored, err := snapshot.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if restored.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", restored.ID, original.ID)
	}
	if restored.Label != original.Label {
		t.Errorf("Label mismatch: got %q, want %q", restored.Label, original.Label)
	}
	if restored.Entries["db/host"] != "localhost" {
		t.Errorf("entry mismatch for db/host")
	}
}

func TestUnmarshal_InvalidJSON(t *testing.T) {
	_, err := snapshot.Unmarshal([]byte(`{invalid}`))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestNew_IDIsDeterministicPerCall(t *testing.T) {
	// Two snapshots with same entries but different timestamps should differ.
	a := snapshot.New("a", map[string]string{"k": "v"})
	b := snapshot.New("b", map[string]string{"k": "v"})
	if a.ID == b.ID {
		t.Error("expected different IDs for snapshots created at different times")
	}
}
