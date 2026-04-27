// Package lint provides validation rules for snapshot state data.
package lint

import (
	"fmt"
	"strings"

	"github.com/user/trackpoint/internal/snapshot"
)

// Severity represents the severity level of a lint violation.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Violation describes a single lint rule failure.
type Violation struct {
	Key      string
	Message  string
	Severity Severity
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Severity, v.Key, v.Message)
}

// Result holds all violations found during linting.
type Result struct {
	Violations []Violation
}

// HasErrors returns true if any violation has error severity.
func (r *Result) HasErrors() bool {
	for _, v := range r.Violations {
		if v.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Run lints a snapshot's state map and returns a Result.
func Run(snap *snapshot.Snapshot) *Result {
	result := &Result{}

	for k, v := range snap.State {
		if strings.TrimSpace(k) == "" {
			result.Violations = append(result.Violations, Violation{
				Key:      "(empty)",
				Message:  "key must not be empty or whitespace",
				Severity: SeverityError,
			})
		}

		if strings.Contains(k, " ") {
			result.Violations = append(result.Violations, Violation{
				Key:      k,
				Message:  "key contains whitespace",
				Severity: SeverityWarning,
			})
		}

		if strings.TrimSpace(v) == "" {
			result.Violations = append(result.Violations, Violation{
				Key:      k,
				Message:  "value is empty or whitespace",
				Severity: SeverityWarning,
			})
		}

		if len(v) > 4096 {
			result.Violations = append(result.Violations, Violation{
				Key:      k,
				Message:  fmt.Sprintf("value exceeds 4096 characters (%d)", len(v)),
				Severity: SeverityWarning,
			})
		}
	}

	return result
}
