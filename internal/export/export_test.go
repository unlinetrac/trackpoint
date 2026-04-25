package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/export"
)

func makeResult() diff.Result {
	return diff.Result{
		Changes: []diff.Change{
			{Key: "cpu", Type: diff.Modified, OldValue: "2", NewValue: "4"},
			{Key: "region", Type: diff.Added, OldValue: nil, NewValue: "us-east-1"},
			{Key: "debug", Type: diff.Removed, OldValue: "true", NewValue: nil},
		},
	}
}

func TestParseFormat_Valid(t *testing.T) {
	for _, tc := range []struct{ in string; want export.Format }{
		{"csv", export.FormatCSV},
		{"CSV", export.FormatCSV},
		{"json", export.FormatJSON},
		{"tsv", export.FormatTSV},
	} {
		got, err := export.ParseFormat(tc.in)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := export.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_CSV(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, makeResult(), export.FormatCSV); err != nil {
		t.Fatalf("Write CSV: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,type,old_value,new_value") {
		t.Errorf("CSV missing header, got:\n%s", out)
	}
	if !strings.Contains(out, "cpu") {
		t.Errorf("CSV missing data row, got:\n%s", out)
	}
}

func TestWrite_TSV(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, makeResult(), export.FormatTSV); err != nil {
		t.Fatalf("Write TSV: %v", err)
	}
	if !strings.Contains(buf.String(), "\t") {
		t.Error("TSV output should contain tab separators")
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := export.Write(&buf, makeResult(), export.FormatJSON); err != nil {
		t.Fatalf("Write JSON: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("JSON output should be an array, got:\n%s", out)
	}
}
