package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/pipeline"
	"github.com/user/trackpoint/internal/snapshot"
)

func init() {
	cmd := &cobra.Command{
		Use:   "pipeline <snapshot-id>",
		Short: "Apply a transformation pipeline to a snapshot and print the result",
		RunE:  runPipeline,
	}

	cmd.Flags().StringSlice("redact", nil, "Comma-separated keys whose values should be redacted")
	cmd.Flags().String("prefix", "", "Keep only keys with this prefix")
	cmd.Flags().Bool("lowercase", false, "Normalize all values to lowercase")
	cmd.Flags().StringSlice("require", nil, "Keys that must be present (error if missing)")

	rootCmd.AddCommand(cmd)
}

func runPipeline(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("snapshot ID required")
	}
	id := args[0]

	store, err := defaultStore()
	if err != nil {
		return err
	}

	snap, err := store.Load(id)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	p := pipeline.New()

	if keys, _ := cmd.Flags().GetStringSlice("redact"); len(keys) > 0 {
		p.Add("redact", pipeline.RedactKeys(keys, "***"))
	}

	if prefix, _ := cmd.Flags().GetString("prefix"); prefix != "" {
		p.Add("filter-prefix", pipeline.FilterKeyPrefix(prefix))
	}

	if lc, _ := cmd.Flags().GetBool("lowercase"); lc {
		p.Add("lowercase", pipeline.NormalizeValues(strings.ToLower))
	}

	if required, _ := cmd.Flags().GetStringSlice("require"); len(required) > 0 {
		p.Add("require-keys", pipeline.RequireKeys(required))
	}

	out, err := p.Run(snap)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "snapshot: %s  label: %s  keys: %d\n", out.ID, out.Label, len(out.Data))
	for _, k := range sortedKeys(out.Data) {
		fmt.Fprintf(os.Stdout, "  %-40s %s\n", k, out.Data[k])
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort — snapshots are small
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// defaultStore is expected to be provided by root.go or snapshot.go.
var _ = snapshot.Snapshot{}
