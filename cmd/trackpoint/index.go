package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"trackpoint/internal/snapshot"
)

func init() {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "Query snapshots using an in-memory label index",
	}

	var labelKey, labelValue, labelKeyOnly string

	lookupCmd := &cobra.Command{
		Use:   "lookup",
		Short: "Find snapshots by label key=value or key",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := defaultStore()
			if err != nil {
				return err
			}
			ids, err := store.List()
			if err != nil {
				return fmt.Errorf("list snapshots: %w", err)
			}
			var snaps []*snapshot.Snapshot
			for _, id := range ids {
				s, err := store.Load(id)
				if err != nil {
					continue
				}
				snaps = append(snaps, s)
			}
			idx := snapshot.NewIndex(snaps)

			var results []*snapshot.Snapshot
			switch {
			case labelKeyOnly != "":
				results = idx.FindByLabelKey(labelKeyOnly)
			case labelKey != "" && labelValue != "":
				results = idx.FindByLabel(labelKey, labelValue)
			default:
				return fmt.Errorf("provide --key or both --label-key and --label-value")
			}

			if len(results) == 0 {
				fmt.Println("no snapshots matched")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tCREATED\tLABELS")
			for _, s := range results {
				labels := ""
				for k, v := range s.Labels {
					if labels != "" {
						labels += ","
					}
					labels += k + "=" + v
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, s.CreatedAt.Format("2006-01-02T15:04:05Z"), labels)
			}
			return w.Flush()
		},
	}

	lookupCmd.Flags().StringVar(&labelKey, "label-key", "", "label key to match")
	lookupCmd.Flags().StringVar(&labelValue, "label-value", "", "label value to match (requires --label-key)")
	lookupCmd.Flags().StringVar(&labelKeyOnly, "key", "", "find all snapshots with this label key (any value)")

	indexCmd.AddCommand(lookupCmd)
	rootCmd.AddCommand(indexCmd)
}
