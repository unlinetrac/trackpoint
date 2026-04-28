package schema

import (
	"testing"
)

func makeSchema(rules ...Rule) *Schema {
	return &Schema{Rules: rules}
}

func TestValidate_NoRules_NoViolations(t *testing.T) {
	s := makeSchema()
	data := map[string]string{"foo": "bar"}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 0 {
		t.Errorf("expected no violations, got %d", len(vs))
	}
}

func TestValidate_ValuePattern_Passes(t *testing.T) {
	s := makeSchema(Rule{
		KeyPattern:  "env",
		ValueRegexp: `^(prod|staging|dev)$`,
	})
	data := map[string]string{"env": "prod"}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 0 {
		t.Errorf("expected no violations, got %v", vs)
	}
}

func TestValidate_ValuePattern_Fails(t *testing.T) {
	s := makeSchema(Rule{
		KeyPattern:  "env",
		ValueRegexp: `^(prod|staging|dev)$`,
	})
	data := map[string]string{"env": "unknown"}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 1 {
		t.Errorf("expected 1 violation, got %d", len(vs))
	}
}

func TestValidate_RequiredKey_Missing(t *testing.T) {
	s := makeSchema(Rule{
		KeyPattern: "version",
		Required:   true,
	})
	data := map[string]string{"env": "prod"}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 1 {
		t.Errorf("expected 1 violation, got %d", len(vs))
	}
}

func TestValidate_RequiredKey_Present(t *testing.T) {
	s := makeSchema(Rule{
		KeyPattern: "version",
		Required:   true,
	})
	data := map[string]string{"version": "1.2.3"}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 0 {
		t.Errorf("expected no violations, got %v", vs)
	}
}

func TestValidate_InvalidKeyPattern_ReturnsError(t *testing.T) {
	s := makeSchema(Rule{KeyPattern: "[invalid"})
	_, err := s.Validate(map[string]string{"foo": "bar"})
	if err == nil {
		t.Error("expected error for invalid key pattern, got nil")
	}
}

func TestValidate_WildcardKeyPattern(t *testing.T) {
	s := makeSchema(Rule{
		KeyPattern:  `service\..*`,
		ValueRegexp: `^[a-z0-9-]+$`,
	})
	data := map[string]string{
		"service.name": "my-service",
		"service.port": "INVALID PORT",
	}
	vs, err := s.Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vs) != 1 {
		t.Errorf("expected 1 violation, got %d: %v", len(vs), vs)
	}
}
