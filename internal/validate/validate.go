// Package validate provides a high-level runner that applies snapshot.Validate
// together with optional schema rules and lint checks, returning a unified
// report of all issues found.
package validate

import (
	"fmt"

	"github.com/yourorg/trackpoint/internal/lint"
	"github.com/yourorg/trackpoint/internal/snapshot"
)

// Issue represents a single problem discovered during validation.
type Issue struct {
	Level   string // "error" or "warning"
	Message string
}

// Result collects every issue found for a snapshot.
type Result struct {
	SnapshotID string
	Issues     []Issue
}

// HasErrors reports whether any issue has level "error".
func (r Result) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Level == "error" {
			return true
		}
	}
	return false
}

// Run validates snap structurally and via lint rules, returning a Result.
func Run(snap *snapshot.Snapshot) (Result, error) {
	if snap == nil {
		return Result{}, fmt.Errorf("validate: snapshot is nil")
	}

	res := Result{SnapshotID: snap.ID}

	// Structural validation.
	if err := snapshot.Validate(snap); err != nil {
		if ve, ok := err.(*snapshot.ValidationError); ok {
			for _, f := range ve.Fields {
				res.Issues = append(res.Issues, Issue{Level: "error", Message: f})
			}
		} else {
			res.Issues = append(res.Issues, Issue{Level: "error", Message: err.Error()})
		}
	}

	// Lint warnings.
	lintViolations := lint.Run(snap)
	for _, v := range lintViolations {
		res.Issues = append(res.Issues, Issue{Level: "warning", Message: v.Message})
	}

	return res, nil
}
