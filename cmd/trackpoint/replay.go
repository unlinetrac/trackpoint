package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/history"
	"github.com/user/trackpoint/internal/replay"
	"github.com/user/trackpoint/internal/snapshot"
)

func init() {
	var onlyChanged bool

	cmd := &cobra.Command{
		Use:   "replay [snapshot-ids...]",
		Short: "Replay diffs across a sequence of snapshots",
		Long: `Replay walks a sequence of snapshots in chronological order and prints
the diff between each consecutive pair. If no IDs are supplied the full
stored timeline is used.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplay(args, onlyChanged)
		},
	}

	cmd.Flags().BoolVar(&onlyChanged, "only-changed", false, "skip frames with no changes")
	rootCmd.AddCommand(cmd)
}

func runReplay(ids []string, onlyChanged bool) error {
	store, err := snapshot.NewStore(defaultSnapshotDir())
	if err != nil {
		return err
	}

	var snaps []*snapshot.Snapshot

	if len(ids) == 0 {
		tl, err := history.NewTimeline(store)
		if err != nil {
			return err
		}
		snaps = tl.Snapshots
	} else {
		for _, id := range ids {
			s, err := store.Load(id)
			if err != nil {
				return fmt.Errorf("loading snapshot %q: %w", id, err)
			}
			snaps = append(snaps, s)
		}
	}

	opts := replay.Options{OnlyChanged: onlyChanged}
	return replay.Run(snaps, opts, func(f replay.Frame) error {
		fmt.Fprintf(os.Stdout, "--- frame %d: %s -> %s ---\n",
			f.Index, f.From.ID[:8], f.To.ID[:8])
		fmt.Fprintln(os.Stdout, diff.SprintResult(f.Result))
		return nil
	})
}
