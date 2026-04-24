package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/diff/format"
	"github.com/user/trackpoint/internal/history"
	"github.com/user/trackpoint/internal/snapshot"
)

var (
	historyStorePath string
	historySortTime  bool
)

var historyCmd = &cobra.Command{
	Use:   "history [snapshot-id...]",
	Short: "Show a diff timeline across multiple snapshots",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().StringVar(&historyStorePath, "store", ".trackpoint", "path to snapshot store directory")
	historyCmd.Flags().BoolVar(&historySortTime, "sort-time", false, "sort snapshots by creation time before diffing")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	store := snapshot.NewStore(historyStorePath)

	tl, err := history.NewTimeline(store, args)
	if err != nil {
		return fmt.Errorf("building timeline: %w", err)
	}

	if historySortTime {
		tl = tl.SortedByTime()
	}

	pairs := tl.Pairs()
	if len(pairs) == 0 {
		fmt.Fprintln(os.Stderr, "no pairs to compare")
		return nil
	}

	for i, pair := range pairs {
		result, err := diff.Compare(pair[0], pair[1])
		if err != nil {
			return fmt.Errorf("comparing pair %d: %w", i+1, err)
		}
		fmt.Printf("=== Step %d: %s → %s ===\n", i+1, pair[0].Label, pair[1].Label)
		fmt.Println(format.SprintResult(result))
	}

	return nil
}
