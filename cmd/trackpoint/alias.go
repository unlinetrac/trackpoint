package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourorg/trackpoint/internal/alias"
)

func init() {
	aliasDir := defaultAliasDir()

	setCmd := &cobra.Command{
		Use:   "set <name> <snapshot-id>",
		Short: "Create or update a named alias for a snapshot ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := alias.NewStore(aliasDir)
			if err := st.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("alias %q -> %s\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Resolve an alias to its snapshot ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := alias.NewStore(aliasDir)
			id, err := st.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Remove a named alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st := alias.NewStore(aliasDir)
			return st.Delete(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all defined aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := alias.NewStore(aliasDir)
			names, err := st.List()
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Println("no aliases defined")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSNAPSHOT ID")
			for _, n := range names {
				id, _ := st.Get(n)
				fmt.Fprintf(w, "%s\t%s\n", n, id)
			}
			return w.Flush()
		},
	}

	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage human-readable aliases for snapshot IDs",
	}
	aliasCmd.AddCommand(setCmd, getCmd, deleteCmd, listCmd)
	rootCmd.AddCommand(aliasCmd)
}

func defaultAliasDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".trackpoint/aliases"
	}
	return home + "/.trackpoint/aliases"
}
