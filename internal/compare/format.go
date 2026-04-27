package compare

import (
	"fmt"
	"strings"

	"github.com/user/trackpoint/internal/diff"
)

// SprintResult formats a compare.Result as a human-readable string.
func SprintResult(r Result) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Range compare: %s → %s (%d snapshots)\n",
		r.FromID, r.ToID, r.SnapshotCount)

	if !r.HasChanges() {
		sb.WriteString("  no net changes detected\n")
		return sb.String()
	}

	for _, c := range r.Diff.Added {
		fmt.Fprintf(&sb, "  + %s = %s\n", c.Key, c.To)
	}
	for _, c := range r.Diff.Removed {
		fmt.Fprintf(&sb, "  - %s (was %s)\n", c.Key, c.From)
	}
	for _, c := range r.Diff.Modified {
		fmt.Fprintf(&sb, "  ~ %s: %s → %s\n", c.Key, c.From, c.To)
	}

	total := len(r.Diff.Added) + len(r.Diff.Removed) + len(r.Diff.Modified)
	fmt.Fprintf(&sb, "  %d change(s) total\n", total)
	return sb.String()
}

// SprintSummary returns a compact one-line summary of a compare.Result.
func SprintSummary(r Result) string {
	if !r.HasChanges() {
		return fmt.Sprintf("%s → %s: no changes", r.FromID, r.ToID)
	}
	return fmt.Sprintf("%s → %s: +%d -%d ~%d",
		r.FromID, r.ToID,
		len(r.Diff.Added),
		len(r.Diff.Removed),
		len(r.Diff.Modified),
	)
}

// ensure diff.Change fields compile — referenced via r.Diff
var _ = diff.Change{}
