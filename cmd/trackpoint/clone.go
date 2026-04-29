package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"trackpoint/internal/snapshot"
)

func init() {
	var overrides []string
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "clone <snapshot-id>",
		Short: "Clone an existing snapshot, optionally overriding key=value pairs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloneSnapshot(args[0], overrides, outputJSON)
		},
	}

	cmd.Flags().StringArrayVarP(&overrides, "set", "s", nil, "Override key=value pairs in the cloned snapshot")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output cloned snapshot as JSON")

	rootCmd.AddCommand(cmd)
}

func runCloneSnapshot(id string, overrides []string, outputJSON bool) error {
	store := snapshot.NewStore(defaultSnapshotDir())

	orig, err := store.Load(id)
	if err != nil {
		return fmt.Errorf("load snapshot %q: %w", id, err)
	}

	var cloned *snapshot.Snapshot
	if len(overrides) > 0 {
		parsed, err := parseKeyValuePairs(overrides)
		if err != nil {
			return fmt.Errorf("parse overrides: %w", err)
		}
		cloned, err = snapshot.CloneWithOverrides(orig, parsed)
		if err != nil {
			return fmt.Errorf("clone with overrides: %w", err)
		}
	} else {
		cloned, err = snapshot.Clone(orig)
		if err != nil {
			return fmt.Errorf("clone snapshot: %w", err)
		}
	}

	if err := store.Save(cloned); err != nil {
		return fmt.Errorf("save cloned snapshot: %w", err)
	}

	if outputJSON {
		return json.NewEncoder(os.Stdout).Encode(cloned)
	}

	fmt.Printf("cloned snapshot saved: %s (from %s)\n", cloned.ID, orig.ID)
	return nil
}
