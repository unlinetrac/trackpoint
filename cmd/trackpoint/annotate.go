package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/trackpoint/internal/annotate"
)

func defaultAnnotateDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".trackpoint/annotations"
	}
	return filepath.Join(home, ".trackpoint", "annotations")
}

func init() {
	var (
		annotateDir string
		snapshotID  string
	)

	// annotate set
	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set an annotation on a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetAnnotation(annotateDir, snapshotID, args[0], args[1])
		},
	}
	setCmd.Flags().StringVar(&annotateDir, "dir", defaultAnnotateDir(), "Directory to store annotations")
	setCmd.Flags().StringVar(&snapshotID, "snapshot", "", "Snapshot ID to annotate (required)")
	_ = setCmd.MarkFlagRequired("snapshot")

	// annotate get
	getCmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get annotations for a snapshot",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := ""
			if len(args) == 1 {
				key = args[0]
			}
			return runGetAnnotation(annotateDir, snapshotID, key)
		},
	}
	getCmd.Flags().StringVar(&annotateDir, "dir", defaultAnnotateDir(), "Directory to store annotations")
	getCmd.Flags().StringVar(&snapshotID, "snapshot", "", "Snapshot ID to look up (required)")
	_ = getCmd.MarkFlagRequired("snapshot")

	// annotate delete
	deleteCmd := &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete an annotation key from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteAnnotation(annotateDir, snapshotID, args[0])
		},
	}
	deleteCmd.Flags().StringVar(&annotateDir, "dir", defaultAnnotateDir(), "Directory to store annotations")
	deleteCmd.Flags().StringVar(&snapshotID, "snapshot", "", "Snapshot ID to modify (required)")
	_ = deleteCmd.MarkFlagRequired("snapshot")

	// annotate parent command
	annotateCmd := &cobra.Command{
		Use:   "annotate",
		Short: "Manage annotations attached to snapshots",
	}
	annotateCmd.AddCommand(setCmd, getCmd, deleteCmd)
	rootCmd.AddCommand(annotateCmd)
}

func runSetAnnotation(dir, snapshotID, key, value string) error {
	store, err := annotate.NewStore(dir)
	if err != nil {
		return fmt.Errorf("open annotation store: %w", err)
	}

	annotations, err := store.Get(snapshotID)
	if err != nil {
		return fmt.Errorf("load annotations: %w", err)
	}

	annotations[key] = value

	if err := store.Set(snapshotID, annotations); err != nil {
		return fmt.Errorf("save annotation: %w", err)
	}

	fmt.Printf("annotation set: %s = %s\n", key, value)
	return nil
}

func runGetAnnotation(dir, snapshotID, key string) error {
	store, err := annotate.NewStore(dir)
	if err != nil {
		return fmt.Errorf("open annotation store: %w", err)
	}

	annotations, err := store.Get(snapshotID)
	if err != nil {
		return fmt.Errorf("load annotations: %w", err)
	}

	if len(annotations) == 0 {
		fmt.Println("no annotations found")
		return nil
	}

	// If a specific key was requested, print just that value.
	if key != "" {
		v, ok := annotations[key]
		if !ok {
			return fmt.Errorf("annotation key %q not found", key)
		}
		fmt.Println(v)
		return nil
	}

	// Print all annotations in a tidy table.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	fmt.Fprintln(w, strings.Repeat("-", 30))
	for k, v := range annotations {
		fmt.Fprintf(w, "%s\t%s\n", k, v)
	}
	return w.Flush()
}

func runDeleteAnnotation(dir, snapshotID, key string) error {
	store, err := annotate.NewStore(dir)
	if err != nil {
		return fmt.Errorf("open annotation store: %w", err)
	}

	if err := store.Delete(snapshotID, key); err != nil {
		return fmt.Errorf("delete annotation: %w", err)
	}

	fmt.Printf("annotation deleted: %s\n", key)
	return nil
}
