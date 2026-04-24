package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var storeDir string

var rootCmd = &cobra.Command{
	Use:   "trackpoint",
	Short: "A lightweight CLI for tracing infrastructure state changes across deploys",
	Long: `trackpoint captures snapshots of infrastructure state and lets you
diff them across deploys to understand what changed and when.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&storeDir,
		"store",
		".trackpoint",
		"Directory used to persist snapshots",
	)
}

func main() {
	Execute()
}
