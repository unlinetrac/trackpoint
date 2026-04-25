package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
	"github.com/user/trackpoint/internal/watch"
)

var (
	watchInterval int
	watchMax      int
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [snapshot-id...]",
		Short: "Watch for changes between snapshots at a given interval",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runWatch,
	}
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 10, "Poll interval in seconds")
	watchCmd.Flags().IntVarP(&watchMax, "max", "n", 1, "Maximum number of poll cycles (0 = infinite)")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	dir := os.Getenv("TRACKPOINT_STORE")
	if dir == "" {
		dir = ".trackpoint/snapshots"
	}

	store := snapshot.NewStore(dir)
	cfg := watch.Config{
		Interval:  time.Duration(watchInterval) * time.Second,
		MaxChecks: watchMax,
	}

	w := watch.New(store, cfg)
	out := make(chan watch.Change, 32)

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Watch(args, out)
	}()

	for change := range out {
		fmt.Printf("[%s] Changes detected:\n", change.At.Format(time.RFC3339))
		for _, c := range change.Result.Changes {
			switch c.Type {
			case diff.Added:
				fmt.Printf("  + %s = %s\n", c.Key, c.To)
			case diff.Removed:
				fmt.Printf("  - %s = %s\n", c.Key, c.From)
			case diff.Modified:
				fmt.Printf("  ~ %s: %s -> %s\n", c.Key, c.From, c.To)
			}
		}
	}

	return <-errCh
}
