package filter_test

import (
	"testing"

	"github.com/trackpoint/internal/diff"
	"github.com/trackpoint/internal/filter"
)

func makeChanges() []diff.Change {
	return []diff.Change{
		{Key: "app.version", Type: diff.Added, To: "1.2.0"},
		{Key: "app.replicas", Type: diff.Modified, From: "2", To: "3"},
		{Key: "db.host", Type: diff.Removed, From: "old-host"},
		{Key: "db.port", Type: diff.Modified, From: "5432", To: "5433"},
		{Key: "cache.ttl", Type: diff.Added, To: "300"},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	changes := makeChanges()
	got, err := filter.Apply(changes, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(changes) {
		t.Errorf("expected %d changes, got %d", len(changes), len(got))
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	got, err := filter.Apply(makeChanges(), filter.Options{KeyPrefix: "db."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 changes, got %d", len(got))
	}
}

func TestApply_TypeFilter(t *testing.T) {
	got, err := filter.Apply(makeChanges(), filter.Options{Types: []diff.ChangeType{diff.Added}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 added changes, got %d", len(got))
	}
}

func TestApply_KeyPattern(t *testing.T) {
	got, err := filter.Apply(makeChanges(), filter.Options{KeyPattern: `^app\.`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 app changes, got %d", len(got))
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := filter.Apply(makeChanges(), filter.Options{KeyPattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestApply_CombinedPrefixAndType(t *testing.T) {
	got, err := filter.Apply(makeChanges(), filter.Options{
		KeyPrefix: "db.",
		Types:     []diff.ChangeType{diff.Modified},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 change, got %d", len(got))
	}
	if got[0].Key != "db.port" {
		t.Errorf("expected db.port, got %s", got[0].Key)
	}
}
