package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/compare"
	"github.com/user/trackpoint/internal/snapshot"
)

func init() {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "compare <id1> <id2> [id3...]",
		Short: "Compare the net change across a range of snapshots",
		Long: `Compare loads two or more snapshots (oldest-first) and reports
the net difference between the first and last entry in the list.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(args, jsonOut)
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output result as JSON")
	rootCmd.AddCommand(cmd)
}

func runCompare(ids []string, jsonOut bool) error {
	store, err := defaultStore()
	if err != nil {
		return err
	}

	snaps := make([]*snapshot.Snapshot, 0, len(ids))
	for _, id := range ids {
		s, err := store.Load(id)
		if err != nil {
			return fmt.Errorf("compare: load %q: %w", id, err)
		}
		snaps = append(snaps, s)
	}

	res, err := compare.Run(snaps)
	if err != nil {
		return err
	}

	if jsonOut {
		out, err := marshalJSON(res)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(out))
		return nil
	}

	fmt.Print(compare.SprintResult(res))
	return nil
}
