package lint_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/trackpoint/internal/lint"
	"github.com/user/trackpoint/internal/snapshot"
)

func makeSnap(state map[string]string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "test-id",
		CreatedAt: time.Now(),
		State:     state,
	}
}

func TestRun_CleanSnapshot_NoViolations(t *testing.T) {
	snap := makeSnap(map[string]string{
		"region":  "us-east-1",
		"version": "1.2.3",
	})
	result := lint.Run(snap)
	if len(result.Violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(result.Violations))
	}
	if result.HasErrors() {
		t.Fatal("expected no errors")
	}
}

func TestRun_EmptyValue_ReturnsWarning(t *testing.T) {
	snap := makeSnap(map[string]string{
		"region": "",
	})
	result := lint.Run(snap)
	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}
	if result.Violations[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning, got %s", result.Violations[0].Severity)
	}
}

func TestRun_KeyWithWhitespace_ReturnsWarning(t *testing.T) {
	snap := makeSnap(map[string]string{
		"my key": "value",
	})
	result := lint.Run(snap)
	found := false
	for _, v := range result.Violations {
		if v.Key == "my key" && v.Severity == lint.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected whitespace-in-key warning")
	}
}

func TestRun_OversizedValue_ReturnsWarning(t *testing.T) {
	snap := makeSnap(map[string]string{
		"blob": strings.Repeat("x", 5000),
	})
	result := lint.Run(snap)
	if len(result.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(result.Violations))
	}
	if result.Violations[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestRun_HasErrors_ReturnsFalseForWarningsOnly(t *testing.T) {
	snap := makeSnap(map[string]string{
		"k": "",
	})
	result := lint.Run(snap)
	if result.HasErrors() {
		t.Error("expected HasErrors to be false for warning-only result")
	}
}

func TestViolation_String_FormatsCorrectly(t *testing.T) {
	v := lint.Violation{Key: "env", Message: "value is empty", Severity: lint.SeverityWarning}
	s := v.String()
	if s != "[warning] env: value is empty" {
		t.Errorf("unexpected format: %s", s)
	}
}
