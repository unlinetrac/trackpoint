package filter

import (
	"regexp"
	"strings"

	"github.com/trackpoint/internal/diff"
)

// Options holds filtering criteria for diff results.
type Options struct {
	KeyPrefix  string
	KeyPattern string
	Types      []diff.ChangeType
}

// Apply filters a slice of diff.Change according to the given Options.
// Returns only the changes that match all specified criteria.
func Apply(changes []diff.Change, opts Options) ([]diff.Change, error) {
	var re *regexp.Regexp
	if opts.KeyPattern != "" {
		var err error
		re, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			return nil, err
		}
	}

	typeSet := make(map[diff.ChangeType]bool, len(opts.Types))
	for _, t := range opts.Types {
		typeSet[t] = true
	}

	var result []diff.Change
	for _, c := range changes {
		if opts.KeyPrefix != "" && !strings.HasPrefix(c.Key, opts.KeyPrefix) {
			continue
		}
		if re != nil && !re.MatchString(c.Key) {
			continue
		}
		if len(typeSet) > 0 && !typeSet[c.Type] {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}
