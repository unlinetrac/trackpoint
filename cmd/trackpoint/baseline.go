package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/baseline"
	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
)

var baselineDir string

func init() {
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage the baseline snapshot used for comparison",
	}

	setCmd := &cobra.Command{
		Use:   "set <snapshot-id>",
		Short: "Set a snapshot as the current baseline",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetBaseline,
	}

	diffCmd := &cobra.Command{
		Use:   "diff <snapshot-id>",
		Short: "Diff a snapshot against the current baseline",
		Args:  cobra.ExactArgs(1),
		RunE:  runDiffBaseline,
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear the current baseline",
		RunE:  runClearBaseline,
	}

	baselineCmd.PersistentFlags().StringVar(&baselineDir, "dir", ".trackpoint", "storage directory")
	baselineCmd.AddCommand(setCmd, diffCmd, clearCmd)
	rootCmd.AddCommand(baselineCmd)
}

func runSetBaseline(cmd *cobra.Command, args []string) error {
	st := snapshot.NewStore(baselineDir)
	bs := baseline.New(st)
	s, err := st.Load(args[0])
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	if err := bs.Set(s); err != nil {
		return fmt.Errorf("set baseline: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Baseline set to snapshot %s\n", args[0])
	return nil
}

func runDiffBaseline(cmd *cobra.Command, args []string) error {
	st := snapshot.NewStore(baselineDir)
	bs := baseline.New(st)
	base, err := bs.Get()
	if err != nil {
		return fmt.Errorf("get baseline: %w", err)
	}
	current, err := st.Load(args[0])
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	result, err := diff.Compare(base, current)
	if err != nil {
		return fmt.Errorf("compare: %w", err)
	}
	fmt.Fprint(os.Stdout, diff.SprintResult(result))
	return nil
}

func runClearBaseline(cmd *cobra.Command, args []string) error {
	st := snapshot.NewStore(baselineDir)
	bs := baseline.New(st)
	if err := bs.Clear(); err != nil {
		return fmt.Errorf("clear baseline: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Baseline cleared")
	return nil
}
