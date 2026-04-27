// Package search provides full-text style search across captured snapshots.
//
// It allows callers to filter snapshots by:
//   - Key substring: match any snapshot containing a state key with the given substring.
//   - Value substring: match any snapshot containing a state value with the given substring.
//   - Tags: match snapshots whose tag map contains all specified key=value pairs.
//
// Example usage:
//
//	results := search.Run(snaps, search.Options{
//		KeyContains:   "db.",
//		ValueContains: "prod",
//		Tags:          map[string]string{"env": "production"},
//	})
//
// Each Result includes the matching Snapshot and the list of keys that
// satisfied the key/value filter criteria.
package search
