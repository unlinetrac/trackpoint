package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/export"
	"github.com/user/trackpoint/internal/snapshot"
)

var (
	exportFormat string
	exportOutput string
)

func init() {
	cmd := &cobra.Command{
		Use:   "export <from-id> <to-id>",
		Short: "Export a diff between two snapshots to CSV, TSV, or JSON",
		Args:  cobra.ExactArgs(2),
		RunE:  runExport,
	}
	cmd.Flags().StringVarP(&exportFormat, "format", "f", "csv", "Output format: csv, tsv, json")
	cmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Write output to file instead of stdout")
	rootCmd.AddCommand(cmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	store := snapshot.NewStore(storeDir)

	fromSnap, err := store.Load(args[0])
	if err != nil {
		return fmt.Errorf("loading snapshot %q: %w", args[0], err)
	}
	toSnap, err := store.Load(args[1])
	if err != nil {
		return fmt.Errorf("loading snapshot %q: %w", args[1], err)
	}

	fmt, err := export.ParseFormat(exportFormat)
	if err != nil {
		return err
	}

	result := diff.Compare(fromSnap, toSnap)

	w := cmd.OutOrStdout()
	if exportOutput != "" {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if err := export.Write(w, result, fmt); err != nil {
		return fmt.Errorf("writing export: %w", err)
	}
	if exportOutput != "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "exported to %s\n", exportOutput)
	}
	return nil
}
