package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"trackpoint/internal/schema"
	"trackpoint/internal/snapshot"
)

func init() {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a snapshot against a JSON schema file",
	}

	var schemaFile string
	var snapshotID string

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a snapshot against schema rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidateSchema(schemaFile, snapshotID)
		},
	}
	validateCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "path to JSON schema file (required)")
	validateCmd.Flags().StringVarP(&snapshotID, "id", "i", "", "snapshot ID to validate (required)")
	_ = validateCmd.MarkFlagRequired("schema")
	_ = validateCmd.MarkFlagRequired("id")

	schemaCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(schemaCmd)
}

func runValidateSchema(schemaFile, snapshotID string) error {
	raw, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var s schema.Schema
	if err := json.Unmarshal(raw, &s); err != nil {
		return fmt.Errorf("parsing schema file: %w", err)
	}

	store := snapshot.NewStore(defaultSnapshotDir())
	snap, err := store.Load(snapshotID)
	if err != nil {
		return fmt.Errorf("loading snapshot %q: %w", snapshotID, err)
	}

	violations, err := s.Validate(snap.Data)
	if err != nil {
		return fmt.Errorf("schema compilation error: %w", err)
	}

	if len(violations) == 0 {
		fmt.Println("✓ snapshot is valid")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema violations (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v)
	}
	return fmt.Errorf("%d violation(s) found", len(violations))
}
