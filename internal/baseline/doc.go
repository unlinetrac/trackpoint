// Package baseline provides management of a pinned "baseline" snapshot.
//
// A baseline is a snapshot that serves as the stable reference point for
// infrastructure state comparisons. Rather than always diffing the most recent
// two snapshots, operators can pin a known-good state as the baseline and
// measure every subsequent snapshot against it.
//
// Usage:
//
//	store := snapshot.NewStore(".trackpoint")
//	bs    := baseline.New(store)
//
//	// Pin a snapshot as the baseline.
//	_ = bs.Set(mySnapshot)
//
//	// Retrieve the baseline for comparison.
//	base, _ := bs.Get()
//
//	// Check whether a baseline has been set.
//	if bs.IsSet() { ... }
//
//	// Remove the baseline.
//	_ = bs.Clear()
package baseline
