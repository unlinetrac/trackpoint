// Package replay walks a chronologically ordered slice of snapshots and
// emits a [Frame] for each consecutive pair, allowing callers to observe
// how infrastructure state evolved over time.
//
// Basic usage:
//
//	err := replay.Run(snapshots, replay.Options{OnlyChanged: true}, func(f replay.Frame) error {
//		fmt.Println(diff.SprintResult(f.Result))
//		return nil
//	})
//
// The Options.OnlyChanged flag causes frames with an empty changeset to be
// skipped, which is useful when replaying long histories that contain many
// identical consecutive snapshots.
package replay
