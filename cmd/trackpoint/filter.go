package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trackpoint/internal/diff"
	"github.com/trackpoint/internal/filter"
	"github.com/trackpoint/internal/report"
	"github.com/trackpoint/internal/snapshot"
)

var (
	filterFrom    string
	filterTo      string
	filterPrefix  string
	filterPattern string
	filterTypes   []string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Diff two snapshots with optional key/type filtering",
	RunE:  runFilter,
}

func init() {
	filterCmd.Flags().StringVar(&filterFrom, "from", "", "ID of the base snapshot (required)")
	filterCmd.Flags().StringVar(&filterTo, "to", "", "ID of the target snapshot (required)")
	filterCmd.Flags().StringVar(&filterPrefix, "prefix", "", "Only include keys with this prefix")
	filterCmd.Flags().StringVar(&filterPattern, "pattern", "", "Only include keys matching this regex")
	filterCmd.Flags().StringSliceVar(&filterTypes, "type", nil, "Change types to include: added,removed,modified")
	_ = filterCmd.MarkFlagRequired("from")
	_ = filterCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(filterCmd)
}

func runFilter(cmd *cobra.Command, args []string) error {
	store := snapshot.NewStore(storeDir)

	fromSnap, err := store.Load(filterFrom)
	if err != nil {
		return fmt.Errorf("loading --from snapshot: %w", err)
	}
	toSnap, err := store.Load(filterTo)
	if err != nil {
		return fmt.Errorf("loading --to snapshot: %w", err)
	}

	result := diff.Compare(fromSnap, toSnap)

	var changeTypes []diff.ChangeType
	for _, t := range filterTypes {
		changeTypes = append(changeTypes, diff.ChangeType(strings.ToLower(t)))
	}

	opts := filter.Options{
		KeyPrefix:  filterPrefix,
		KeyPattern: filterPattern,
		Types:      changeTypes,
	}

	filtered, err := filter.Apply(result.Changes, opts)
	if err != nil {
		return fmt.Errorf("applying filter: %w", err)
	}
	result.Changes = filtered

	rep, err := report.New(fromSnap, toSnap)
	if err != nil {
		return err
	}
	rep.Result = result

	fmt.Print(rep.Sprint())
	return nil
}
