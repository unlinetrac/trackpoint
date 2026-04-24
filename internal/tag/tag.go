package tag

import (
	"errors"
	"regexp"
	"strings"
)

var validTagKey = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// Tags is a map of string key-value labels attached to a snapshot.
type Tags map[string]string

// Parse parses a slice of "key=value" strings into a Tags map.
func Parse(pairs []string) (Tags, error) {
	tags := make(Tags, len(pairs))
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid tag format: expected key=value, got: " + pair)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if !validTagKey.MatchString(key) {
			return nil, errors.New("invalid tag key: " + key)
		}
		tags[key] = val
	}
	return tags, nil
}

// Match reports whether the given tags satisfy all of the filter tags.
// Every key in filter must exist in tags with the same value.
func Match(tags Tags, filter Tags) bool {
	for k, v := range filter {
		if tags[k] != v {
			return false
		}
	}
	return true
}

// Merge returns a new Tags map combining base and override.
// Keys in override take precedence.
func Merge(base, override Tags) Tags {
	out := make(Tags, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}
