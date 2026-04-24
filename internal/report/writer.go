package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/trackpoint/internal/diff"
)

func writeJSON(w io.Writer, r *Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		return fmt.Errorf("report: failed to encode JSON: %w", err)
	}
	return nil
}

func writeText(w io.Writer, r *Report) error {
	fmt.Fprintf(w, "Snapshot diff\n")
	fmt.Fprintf(w, "  from : %s\n", r.From)
	fmt.Fprintf(w, "  to   : %s\n", r.To)
	fmt.Fprintf(w, "  at   : %s\n\n", r.Timestamp.Format("2006-01-02 15:04:05 UTC"))

	if !r.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return nil
	}

	fmt.Fprintf(w, "Changes (%d):\n", len(r.Result.Changes))
	for _, c := range r.Result.Changes {
		switch c.Type {
		case diff.ChangeAdded:
			fmt.Fprintf(w, "  + %s = %v\n", c.Key, c.NewValue)
		case diff.ChangeRemoved:
			fmt.Fprintf(w, "  - %s = %v\n", c.Key, c.OldValue)
		case diff.ChangeModified:
			fmt.Fprintf(w, "  ~ %s: %v -> %v\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return nil
}
