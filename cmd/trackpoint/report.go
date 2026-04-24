package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trackpoint/internal/report"
	"github.com/trackpoint/internal/snapshot"
)

var reportFormat string

var reportCmd = &cobra.Command{
	Use:   "report <from-id> <to-id>",
	Short: "Generate a diff report between two snapshots",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := snapshot.NewStore(storeDir)
		if err != nil {
			return fmt.Errorf("failed to open store: %w", err)
		}

		from, err := store.Load(args[0])
		if err != nil {
			return fmt.Errorf("failed to load snapshot %q: %w", args[0], err)
		}

		to, err := store.Load(args[1])
		if err != nil {
			return fmt.Errorf("failed to load snapshot %q: %w", args[1], err)
		}

		r, err := report.New(from, to)
		if err != nil {
			return fmt.Errorf("failed to create report: %w", err)
		}

		if err := r.Write(os.Stdout, report.Format(reportFormat)); err != nil {
			return fmt.Errorf("failed to write report: %w", err)
		}

		if r.HasChanges() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "text", "Output format: text or json")
	rootCmd.AddCommand(reportCmd)
}
