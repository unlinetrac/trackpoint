package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/snapshot"
)

func init() {
	var labelFilter string
	var since string
	var until string
	var tags []string
	var limit int

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query snapshots by label, time range, or tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := defaultSnapshotDir()
			store, err := snapshot.NewStore(dir)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			opts := snapshot.QueryOptions{
				Label: labelFilter,
				Limit: limit,
			}

			if since != "" {
				t, err := time.Parse(time.RFC3339, since)
				if err != nil {
					return fmt.Errorf("invalid --since: %w", err)
				}
				opts.Since = &t
			}
			if until != "" {
				t, err := time.Parse(time.RFC3339, until)
				if err != nil {
					return fmt.Errorf("invalid --until: %w", err)
				}
				opts.Until = &t
			}
			if len(tags) > 0 {
				parsed, err := parseKeyValuePairs(tags)
				if err != nil {
					return fmt.Errorf("invalid --tag: %w", err)
				}
				opts.Tags = parsed
			}

			results, err := snapshot.Query(store, opts)
			if err != nil {
				return fmt.Errorf("query: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(os.Stderr, "no snapshots matched")
				return nil
			}
			for _, s := range results {
				fmt.Printf("%s\t%s\t%s\n", s.ID, s.Label, s.CreatedAt.Format(time.RFC3339))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&labelFilter, "label", "", "filter by label substring")
	cmd.Flags().StringVar(&since, "since", "", "include snapshots at or after this RFC3339 time")
	cmd.Flags().StringVar(&until, "until", "", "include snapshots at or before this RFC3339 time")
	cmd.Flags().StringArrayVar(&tags, "tag", nil, "filter by tag key=value (repeatable)")
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum number of results (0 = unlimited)")

	rootCmd.AddCommand(cmd)
}
