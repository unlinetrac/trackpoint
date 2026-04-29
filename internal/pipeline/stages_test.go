package pipeline_test

import (
	"strings"
	"testing"

	"github.com/user/trackpoint/internal/pipeline"
	"github.com/user/trackpoint/internal/snapshot"
)

func TestRedactKeys_ReplacesMatchingValues(t *testing.T) {
	snap := makeSnap("s", map[string]string{"secret": "pass123", "host": "localhost"})
	fn := pipeline.RedactKeys([]string{"secret"}, "***")
	out, err := fn(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data["secret"] != "***" {
		t.Errorf("expected redacted value, got %q", out.Data["secret"])
	}
	if out.Data["host"] != "localhost" {
		t.Errorf("expected host unchanged, got %q", out.Data["host"])
	}
}

func TestRedactKeys_MissingKeyIsIgnored(t *testing.T) {
	snap := makeSnap("s", map[string]string{"a": "1"})
	fn := pipeline.RedactKeys([]string{"missing"}, "***")
	out, err := fn(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Data["missing"]; ok {
		t.Error("missing key should not be added")
	}
}

func TestFilterKeyPrefix_KeepsMatchingKeys(t *testing.T) {
	snap := makeSnap("s", map[string]string{"app.port": "8080", "db.host": "localhost", "app.name": "tp"})
	fn := pipeline.FilterKeyPrefix("app.")
	out, err := fn(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Data) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(out.Data), out.Data)
	}
	if _, ok := out.Data["db.host"]; ok {
		t.Error("db.host should have been filtered out")
	}
}

func TestNormalizeValues_AppliesTransform(t *testing.T) {
	snap := makeSnap("s", map[string]string{"env": "Production", "mode": "DEBUG"})
	fn := pipeline.NormalizeValues(strings.ToLower)
	out, err := fn(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Data["env"] != "production" {
		t.Errorf("expected lowercase, got %q", out.Data["env"])
	}
	if out.Data["mode"] != "debug" {
		t.Errorf("expected lowercase, got %q", out.Data["mode"])
	}
}

func TestRequireKeys_AllPresent_NoError(t *testing.T) {
	snap := makeSnap("s", map[string]string{"host": "localhost", "port": "8080"})
	fn := pipeline.RequireKeys([]string{"host", "port"})
	_, err := fn(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireKeys_MissingKey_ReturnsError(t *testing.T) {
	snap := makeSnap("s", map[string]string{"host": "localhost"})
	fn := pipeline.RequireKeys([]string{"host", "port"})
	_, err := fn(snap)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "port") {
		t.Errorf("expected key name in error, got %v", err)
	}
}

func TestPipeline_RedactThenFilter_Composes(t *testing.T) {
	snap := makeSnap("s", map[string]string{
		"app.secret": "topsecret",
		"app.host":   "localhost",
		"db.pass":    "dbpass",
	})
	p := pipeline.New().
		Add("redact", pipeline.RedactKeys([]string{"app.secret", "db.pass"}, "***")).
		Add("filter", pipeline.FilterKeyPrefix("app."))
	out, err := p.Run(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.Data["db.pass"]; ok {
		t.Error("db.pass should have been filtered")
	}
	if out.Data["app.secret"] != "***" {
		t.Errorf("expected redacted, got %q", out.Data["app.secret"])
	}
	if out.Data["app.host"] != "localhost" {
		t.Errorf("expected host unchanged, got %q", out.Data["app.host"])
	}
}
