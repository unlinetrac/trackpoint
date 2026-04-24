package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"trackpoint/internal/tag"
)

var tagDir string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on snapshots",
	}

	setCmd := &cobra.Command{
		Use:   "set <snapshot-id> [key=value ...]",
		Short: "Attach tags to a snapshot",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runSetTags,
	}
	setCmd.Flags().StringVar(&tagDir, "dir", ".trackpoint", "directory where snapshots are stored")

	getCmd := &cobra.Command{
		Use:   "get <snapshot-id>",
		Short: "Print tags for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE:  runGetTags,
	}
	getCmd.Flags().StringVar(&tagDir, "dir", ".trackpoint", "directory where snapshots are stored")

	tagCmd.AddCommand(setCmd, getCmd)
	rootCmd.AddCommand(tagCmd)
}

func runSetTags(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]
	pairs := args[1:]

	newTags, err := tag.Parse(pairs)
	if err != nil {
		return fmt.Errorf("parsing tags: %w", err)
	}

	store := tag.NewTagStore(tagDir)
	existing, err := store.Load(snapshotID)
	if err != nil {
		return fmt.Errorf("loading existing tags: %w", err)
	}

	merged := tag.Merge(existing, newTags)
	if err := store.Save(snapshotID, merged); err != nil {
		return fmt.Errorf("saving tags: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Tags updated for snapshot %s\n", snapshotID)
	return nil
}

func runGetTags(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]
	store := tag.NewTagStore(tagDir)

	tags, err := store.Load(snapshotID)
	if err != nil {
		return fmt.Errorf("loading tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Fprintf(os.Stdout, "No tags for snapshot %s\n", snapshotID)
		return nil
	}

	for k, v := range tags {
		fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
	}
	return nil
}
