package filter

import (
	"fmt"
	"strings"

	"github.com/trackpoint/internal/diff"
)

// ParseTypes converts a slice of raw strings into validated ChangeType values.
// Returns an error if any string is not a recognised change type.
func ParseTypes(raw []string) ([]diff.ChangeType, error) {
	valid := map[string]diff.ChangeType{
		"added":    diff.Added,
		"removed":  diff.Removed,
		"modified": diff.Modified,
	}

	out := make([]diff.ChangeType, 0, len(raw))
	for _, s := range raw {
		norm := strings.ToLower(strings.TrimSpace(s))
		ct, ok := valid[norm]
		if !ok {
			return nil, fmt.Errorf("unknown change type %q: must be one of added, removed, modified", s)
		}
		out = append(out, ct)
	}
	return out, nil
}

// Summary returns a human-readable description of the active filter options.
func (o Options) Summary() string {
	var parts []string
	if o.KeyPrefix != "" {
		parts = append(parts, fmt.Sprintf("prefix=%q", o.KeyPrefix))
	}
	if o.KeyPattern != "" {
		parts = append(parts, fmt.Sprintf("pattern=%q", o.KeyPattern))
	}
	if len(o.Types) > 0 {
		ts := make([]string, len(o.Types))
		for i, t := range o.Types {
			ts[i] = string(t)
		}
		parts = append(parts, fmt.Sprintf("types=[%s]", strings.Join(ts, ",")))
	}
	if len(parts) == 0 {
		return "no filters applied"
	}
	return strings.Join(parts, " ")
}
