package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trackpoint/internal/search"
	"github.com/trackpoint/internal/snapshot"
)

var (
	searchKeyContains   string
	searchValueContains string
	searchTags          []string
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search snapshots by key, value, or tag",
		RunE:  runSearch,
	}

	searchCmd.Flags().StringVar(&searchKeyContains, "key", "", "Filter by key substring")
	searchCmd.Flags().StringVar(&searchValueContains, "value", "", "Filter by value substring")
	searchCmd.Flags().StringArrayVar(&searchTags, "tag", nil, "Filter by tag (key=value)")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	store := snapshot.NewStore(defaultSnapshotDir())

	ids, err := store.List()
	if err != nil {
		return fmt.Errorf("listing snapshots: %w", err)
	}

	var snaps []*snapshot.Snapshot
	for _, id := range ids {
		snap, err := store.Load(id)
		if err != nil {
			continue
		}
		snaps = append(snaps, snap)
	}

	tagMap := make(map[string]string)
	for _, t := range searchTags {
		parts := strings.SplitN(t, "=", 2)
		if len(parts) == 2 {
			tagMap[parts[0]] = parts[1]
		}
	}

	opts := search.Options{
		KeyContains:   searchKeyContains,
		ValueContains: searchValueContains,
		Tags:          tagMap,
	}

	results := search.Run(snaps, opts)
	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matching snapshots found")
		return nil
	}

	for _, r := range results {
		fmt.Printf("snapshot: %s  created: %s  matched_keys: %s\n",
			r.Snapshot.ID,
			r.Snapshot.CreatedAt.Format("2006-01-02T15:04:05Z"),
			strings.Join(r.MatchedKeys, ", "),
		)
	}
	return nil
}
