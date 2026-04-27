package search

import (
	"strings"

	"github.com/trackpoint/internal/snapshot"
)

// Options holds the search parameters.
type Options struct {
	KeyContains   string
	ValueContains string
	Tags          map[string]string
}

// Result holds a matching snapshot and the keys that matched.
type Result struct {
	Snapshot  *snapshot.Snapshot
	MatchedKeys []string
}

// Run searches through a list of snapshots and returns results
// that satisfy all provided options.
func Run(snaps []*snapshot.Snapshot, opts Options) []Result {
	var results []Result

	for _, snap := range snaps {
		if snap == nil {
			continue
		}

		if !tagsMatch(snap.Tags, opts.Tags) {
			continue
		}

		matched := matchedKeys(snap.State, opts.KeyContains, opts.ValueContains)
		if len(matched) == 0 && (opts.KeyContains != "" || opts.ValueContains != "") {
			continue
		}

		results = append(results, Result{
			Snapshot:    snap,
			MatchedKeys: matched,
		})
	}

	return results
}

func matchedKeys(state map[string]string, keyContains, valueContains string) []string {
	var keys []string
	for k, v := range state {
		if keyContains != "" && !strings.Contains(k, keyContains) {
			continue
		}
		if valueContains != "" && !strings.Contains(v, valueContains) {
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func tagsMatch(snapTags, filterTags map[string]string) bool {
	for k, v := range filterTags {
		if snapTags[k] != v {
			return false
		}
	}
	return true
}
