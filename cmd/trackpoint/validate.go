package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/trackpoint/internal/snapshot"
	"github.com/yourorg/trackpoint/internal/validate"
)

func init() {
	var storeDir string

	cmd := &cobra.Command{
		Use:   "validate <snapshot-id>",
		Short: "Validate a snapshot for structural and lint issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(storeDir, args[0])
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", defaultStoreDir(), "path to snapshot store")
	rootCmd.AddCommand(cmd)
}

func runValidate(storeDir, id string) error {
	store, err := snapshot.NewStore(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	snap, err := store.Load(id)
	if err != nil {
		return fmt.Errorf("load snapshot %q: %w", id, err)
	}

	res, err := validate.Run(snap)
	if err != nil {
		return err
	}

	if len(res.Issues) == 0 {
		fmt.Printf("snapshot %s: OK\n", id)
		return nil
	}

	for _, issue := range res.Issues {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", issue.Level, issue.Message)
	}

	if res.HasErrors() {
		return fmt.Errorf("snapshot %s has validation errors", id)
	}

	// Warnings only — exit cleanly but still print them.
	return nil
}
