package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/user/trackpoint/internal/lock"
)

func defaultLockDir() string {
	return filepath.Join(defaultDataDir(), "locks")
}

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage snapshot locks",
	}

	ackCmd := &cobra.Command{
		Use:   "acquire <snapshot-id>",
		Short: "Acquire a lock on a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAcquireLock(args[0])
		},
	}

	releaseCmd := &cobra.Command{
		Use:   "release <snapshot-id>",
		Short: "Release a lock on a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReleaseLock(args[0])
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <snapshot-id>",
		Short: "Check whether a snapshot is locked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockStatus(args[0])
		},
	}

	lockCmd.AddCommand(ackCmd, releaseCmd, statusCmd)
	rootCmd.AddCommand(lockCmd)
}

func runAcquireLock(id string) error {
	s := lock.NewStore(defaultLockDir())
	if err := s.Lock(id); err != nil {
		return fmt.Errorf("acquire lock: %w", err)
	}
	fmt.Fprintf(os.Stdout, "locked: %s\n", id)
	return nil
}

func runReleaseLock(id string) error {
	s := lock.NewStore(defaultLockDir())
	if err := s.Unlock(id); err != nil {
		return fmt.Errorf("release lock: %w", err)
	}
	fmt.Fprintf(os.Stdout, "unlocked: %s\n", id)
	return nil
}

func runLockStatus(id string) error {
	s := lock.NewStore(defaultLockDir())
	if s.IsLocked(id) {
		fmt.Fprintf(os.Stdout, "%s: locked\n", id)
	} else {
		fmt.Fprintf(os.Stdout, "%s: unlocked\n", id)
	}
	return nil
}
