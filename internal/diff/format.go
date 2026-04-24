package diff

import (
	"fmt"
	"io"
	"strings"
)

// Formatter defines how a diff Result is rendered.
type Formatter interface {
	Format(w io.Writer, r *Result) error
}

// TextFormatter renders diffs as human-readable plain text.
type TextFormatter struct{}

// Format writes a plain-text representation of the diff to w.
func (f *TextFormatter) Format(w io.Writer, r *Result) error {
	fmt.Fprintf(w, "Diff: %s → %s\n", r.FromID, r.ToID)
	if !r.HasChanges() {
		fmt.Fprintln(w, "  (no changes)")
		return nil
	}
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s: %v\n", c.Key, c.NewVal)
		case Removed:
			fmt.Fprintf(w, "  - %s: %v\n", c.Key, c.OldVal)
		case Modified:
			fmt.Fprintf(w, "  ~ %s: %v → %v\n", c.Key, c.OldVal, c.NewVal)
		}
	}
	return nil
}

// SprintResult returns the formatted diff as a string using TextFormatter.
func SprintResult(r *Result) string {
	var sb strings.Builder
	f := &TextFormatter{}
	_ = f.Format(&sb, r)
	return sb.String()
}
