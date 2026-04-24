package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/trackpoint/internal/report"
	"github.com/trackpoint/internal/snapshot"
)

func makeSnap(t *testing.T, label string, state map[string]interface{}) *snapshot.Snapshot {
	t.Helper()
	s, err := snapshot.New(label, state)
	if err != nil {
		t.Fatalf("makeSnap: %v", err)
	}
	return s
}

func TestNew_ReturnsReport(t *testing.T) {
	from := makeSnap(t, "v1", map[string]interface{}{"replicas": 2})
	to := makeSnap(t, "v2", map[string]interface{}{"replicas": 3})

	r, err := report.New(from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.From != from.ID {
		t.Errorf("expected From=%s, got %s", from.ID, r.From)
	}
	if r.To != to.ID {
		t.Errorf("expected To=%s, got %s", to.ID, r.To)
	}
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNew_NilFromReturnsError(t *testing.T) {
	to := makeSnap(t, "v2", map[string]interface{}{})
	_, err := report.New(nil, to)
	if err == nil {
		t.Error("expected error for nil from snapshot")
	}
}

func TestNew_NilToReturnsError(t *testing.T) {
	from := makeSnap(t, "v1", map[string]interface{}{})
	_, err := report.New(from, nil)
	if err == nil {
		t.Error("expected error for nil to snapshot")
	}
}

func TestHasChanges(t *testing.T) {
	from := makeSnap(t, "v1", map[string]interface{}{"x": 1})
	to := makeSnap(t, "v2", map[string]interface{}{"x": 2})

	r, _ := report.New(from, to)
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestHasChanges_NoChanges(t *testing.T) {
	state := map[string]interface{}{"x": 1}
	from := makeSnap(t, "v1", state)
	to := makeSnap(t, "v2", state)

	r, _ := report.New(from, to)
	if r.HasChanges() {
		t.Error("expected HasChanges to be false when state is identical")
	}
}

func TestWrite_TextFormat_NoChanges(t *testing.T) {
	state := map[string]interface{}{"key": "val"}
	from := makeSnap(t, "v1", state)
	to := makeSnap(t, "v2", state)

	r, _ := report.New(from, to)
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatText); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	_ = time.Now()
	from := makeSnap(t, "v1", map[string]interface{}{"a": 1})
	to := makeSnap(t, "v2", map[string]interface{}{"a": 2})

	r, _ := report.New(from, to)
	var buf bytes.Buffer
	if err := r.Write(&buf, report.FormatJSON); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if !strings.Contains(buf.String(), `"from"`) {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	from := makeSnap(t, "v1", map[string]interface{}{})
	to := makeSnap(t, "v2", map[string]interface{}{})

	r, _ := report.New(from, to)
	var buf bytes.Buffer
	if err := r.Write(&buf, report.Format("xml")); err == nil {
		t.Error("expected error for unsupported format")
	}
}
