// Package snapshot provides primitives for capturing, storing, and querying
// point-in-time infrastructure state.
//
// # Query
//
// The Query function allows callers to retrieve a filtered subset of snapshots
// from any store that implements the Loader interface:
//
//	results, err := snapshot.Query(store, snapshot.QueryOptions{
//		Label: "deploy",
//		Tags:  map[string]string{"env": "prod"},
//		Limit: 10,
//	})
//
// Supported filters:
//
//   - Label: case-insensitive substring match on the snapshot label.
//   - Since / Until: time-range bounds using *time.Time pointers.
//   - Tags: all specified key=value pairs must be present on the snapshot.
//   - Limit: cap the number of returned results (0 means no cap).
//
// Filters are AND-ed together. Results are returned in store list order,
// which is lexicographically sorted by snapshot ID.
package snapshot
