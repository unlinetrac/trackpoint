package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/snapshot"
)

var snapshotSource string
var snapshotStoreDir string

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage infrastructure state snapshots",
	}

	takeCmd := &cobra.Command{
		Use:   "take",
		Short: "Take a new snapshot from KEY=VALUE pairs on stdin or flags",
		RunE:  runTakeSnapshot,
	}
	takeCmd.Flags().StringVarP(&snapshotSource, "source", "s", "", "Source label for the snapshot (required)")
	takeCmd.Flags().StringVarP(&snapshotStoreDir, "dir", "d", ".trackpoint", "Directory to store snapshots")
	_ = takeCmd.MarkFlagRequired("source")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List stored snapshot IDs",
		RunE:  runListSnapshots,
	}
	listCmd.Flags().StringVarP(&snapshotStoreDir, "dir", "d", ".trackpoint", "Directory to read snapshots from")

	snapshotCmd.AddCommand(takeCmd, listCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runTakeSnapshot(cmd *cobra.Command, args []string) error {
	state := make(map[string]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid KEY=VALUE pair: %q", arg)
		}
		state[parts[0]] = parts[1]
	}

	snap := snapshot.New(snapshotSource, state)
	store := snapshot.NewStore(snapshotStoreDir)
	if err := store.Save(snap); err != nil {
		return fmt.Errorf("saving snapshot: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "snapshot saved: %s\n", snap.ID)
	return nil
}

func runListSnapshots(cmd *cobra.Command, args []string) error {
	store := snapshot.NewStore(snapshotStoreDir)
	ids, err := store.List()
	if err != nil {
		return fmt.Errorf("listing snapshots: %w", err)
	}
	if len(ids) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no snapshots found")
		return nil
	}
	for _, id := range ids {
		fmt.Fprintln(cmd.OutOrStdout(), id)
	}
	return nil
}
