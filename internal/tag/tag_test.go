package tag

import (
	"testing"
)

func TestParse_ValidPairs(t *testing.T) {
	tags, err := Parse([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", tags["env"])
	}
	if tags["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %s", tags["region"])
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]string{"badtag"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParse_InvalidKey(t *testing.T) {
	_, err := Parse([]string{"bad key=value"})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestParse_Empty(t *testing.T) {
	tags, err := Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestMatch_AllMatch(t *testing.T) {
	tags := Tags{"env": "prod", "region": "us-east-1"}
	filter := Tags{"env": "prod"}
	if !Match(tags, filter) {
		t.Error("expected match")
	}
}

func TestMatch_NoMatch(t *testing.T) {
	tags := Tags{"env": "staging"}
	filter := Tags{"env": "prod"}
	if Match(tags, filter) {
		t.Error("expected no match")
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	tags := Tags{"env": "prod"}
	if !Match(tags, Tags{}) {
		t.Error("empty filter should match everything")
	}
}

func TestMerge_OverrideWins(t *testing.T) {
	base := Tags{"env": "staging", "app": "web"}
	override := Tags{"env": "prod"}
	out := Merge(base, override)
	if out["env"] != "prod" {
		t.Errorf("expected prod, got %s", out["env"])
	}
	if out["app"] != "web" {
		t.Errorf("expected web, got %s", out["app"])
	}
}
