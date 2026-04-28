// Package schema provides key schema validation for snapshots,
// ensuring snapshot data conforms to expected key patterns and value types.
package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for a key pattern.
type Rule struct {
	KeyPattern  string         `json:"key_pattern"`
	Required    bool           `json:"required"`
	ValueRegexp string         `json:"value_regexp"`
	Description string         `json:"description"`
	compiledKey *regexp.Regexp
	compiledVal *regexp.Regexp
}

// Schema holds a set of validation rules.
type Schema struct {
	Rules []Rule `json:"rules"`
}

// Violation represents a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

// Validate checks the given key-value data against the schema rules.
// It returns a list of violations (empty means valid).
func (s *Schema) Validate(data map[string]string) ([]Violation, error) {
	if err := s.compile(); err != nil {
		return nil, err
	}

	var violations []Violation

	for i := range s.Rules {
		rule := &s.Rules[i]
		matched := false

		for key, val := range data {
			if !rule.compiledKey.MatchString(key) {
				continue
			}
			matched = true
			if rule.compiledVal != nil && !rule.compiledVal.MatchString(val) {
				violations = append(violations, Violation{
					Key:     key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.ValueRegexp),
				})
			}
		}

		if rule.Required && !matched {
			violations = append(violations, Violation{
				Key:     rule.KeyPattern,
				Message: fmt.Sprintf("required key pattern %q not found in snapshot", rule.KeyPattern),
			})
		}
	}

	return violations, nil
}

func (s *Schema) compile() error {
	for i := range s.Rules {
		r := &s.Rules[i]
		if r.compiledKey != nil {
			continue
		}
		pattern := r.KeyPattern
		if !strings.HasPrefix(pattern, "^") {
			pattern = "^" + pattern
		}
		if !strings.HasSuffix(pattern, "$") {
			pattern = pattern + "$"
		}
		ck, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid key pattern %q: %w", r.KeyPattern, err)
		}
		r.compiledKey = ck
		if r.ValueRegexp != "" {
			cv, err := regexp.Compile(r.ValueRegexp)
			if err != nil {
				return fmt.Errorf("invalid value regexp %q: %w", r.ValueRegexp, err)
			}
			r.compiledVal = cv
		}
	}
	return nil
}
