// Package replay provides functionality for replaying snapshot diffs
// across a sequence of snapshots in chronological order.
package replay

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
)

// Frame represents a single step in a replay sequence.
type Frame struct {
	Index  int
	From   *snapshot.Snapshot
	To     *snapshot.Snapshot
	Result diff.Result
}

// Options controls replay behaviour.
type Options struct {
	// OnlyChanged skips frames where no changes were detected.
	OnlyChanged bool
}

// Run iterates over the provided snapshots in order and emits a Frame for
// each consecutive pair. The snapshots slice must contain at least two
// entries. Frames are delivered to the provided callback; returning an
// error from the callback halts the replay.
func Run(snapshots []*snapshot.Snapshot, opts Options, fn func(Frame) error) error {
	if len(snapshots) < 2 {
		return errors.New("replay: at least two snapshots are required")
	}

	for i := 1; i < len(snapshots); i++ {
		from := snapshots[i-1]
		to := snapshots[i]

		result, err := diff.Compare(from, to)
		if err != nil {
			return fmt.Errorf("replay: comparing snapshots at index %d: %w", i, err)
		}

		if opts.OnlyChanged && len(result.Changes) == 0 {
			continue
		}

		frame := Frame{
			Index:  i,
			From:   from,
			To:     to,
			Result: result,
		}

		if err := fn(frame); err != nil {
			return err
		}
	}

	return nil
}
