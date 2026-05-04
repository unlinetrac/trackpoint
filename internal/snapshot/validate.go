package snapshot

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError holds a list of field-level issues found in a Snapshot.
type ValidationError struct {
	Fields []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("snapshot validation failed: %s", strings.Join(e.Fields, "; "))
}

// Validate checks that a Snapshot is structurally sound and returns a
// *ValidationError listing every problem found, or nil if the snapshot is valid.
func Validate(s *Snapshot) error {
	if s == nil {
		return errors.New("snapshot is nil")
	}

	var problems []string

	if strings.TrimSpace(s.ID) == "" {
		problems = append(problems, "id must not be empty")
	}

	if s.CreatedAt.IsZero() {
		problems = append(problems, "created_at must not be zero")
	}

	if s.Data == nil {
		problems = append(problems, "data must not be nil")
	}

	for k, v := range s.Data {
		if strings.TrimSpace(k) == "" {
			problems = append(problems, "data contains a blank key")
		}
		if len(v) > 4096 {
			problems = append(problems, fmt.Sprintf("value for key %q exceeds 4096 bytes", k))
		}
	}

	for k, v := range s.Labels {
		if strings.TrimSpace(k) == "" {
			problems = append(problems, "labels contain a blank key")
		}
		if strings.TrimSpace(v) == "" {
			problems = append(problems, fmt.Sprintf("label %q has an empty value", k))
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Fields: problems}
	}
	return nil
}
