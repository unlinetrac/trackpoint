package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/rollback"
	"github.com/user/trackpoint/internal/snapshot"
)

func init() {
	var noDiff bool

	cmd := &cobra.Command{
		Use:   "rollback <snapshot-id>",
		Short: "Find the best rollback target relative to a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRollback(args[0], noDiff)
		},
	}

	cmd.Flags().BoolVar(&noDiff, "no-diff", false, "only suggest a target with no state changes vs current")
	rootCmd.AddCommand(cmd)
}

func runRollback(currentID string, noDiff bool) error {
	store, err := snapshot.NewStore(defaultSnapshotDir())
	if err != nil {
		return fmt.Errorf("open snapshot store: %w", err)
	}

	res, err := rollback.Find(store, currentID, noDiff)
	if err != nil {
		return fmt.Errorf("rollback: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Rollback target: %s (created %s)\n",
		res.Target.ID, res.Target.CreatedAt.Format("2006-01-02 15:04:05"))

	if len(res.Changes) == 0 {
		fmt.Fprintln(os.Stdout, "No state changes relative to current snapshot.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "\n%d change(s) relative to current snapshot:\n", len(res.Changes))
	for _, c := range res.Changes {
		switch c.Type {
		case diff.Added:
			fmt.Fprintf(os.Stdout, "  + %s = %s\n", c.Key, c.NewValue)
		case diff.Removed:
			fmt.Fprintf(os.Stdout, "  - %s (was %s)\n", c.Key, c.OldValue)
		case diff.Modified:
			fmt.Fprintf(os.Stdout, "  ~ %s: %s → %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return nil
}
