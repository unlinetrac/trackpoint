// Package snapshot provides types and utilities for capturing and managing
// infrastructure state snapshots.
package snapshot

import (
	"errors"
	"fmt"
)

// PatchOp represents a single patch operation to apply to a snapshot's data.
type PatchOp struct {
	// Key is the data key to operate on.
	Key string
	// Op is the operation type: "set", "delete", or "rename".
	Op string
	// Value is the new value for "set" operations.
	Value string
	// NewKey is the destination key name for "rename" operations.
	NewKey string
}

// Patch applies a slice of PatchOps to the given snapshot and returns a new
// snapshot with the modifications applied. The original snapshot is not mutated.
func Patch(s *Snapshot, ops []PatchOp) (*Snapshot, error) {
	if s == nil {
		return nil, errors.New("patch: snapshot must not be nil")
	}
	if len(ops) == 0 {
		return Clone(s)
	}

	// Start from a deep copy so the original is never mutated.
	patched, err := Clone(s)
	if err != nil {
		return nil, fmt.Errorf("patch: clone failed: %w", err)
	}

	for _, op := range ops {
		if op.Key == "" {
			return nil, fmt.Errorf("patch: op %q has empty key", op.Op)
		}
		switch op.Op {
		case "set":
			patched.Data[op.Key] = op.Value
		case "delete":
			delete(patched.Data, op.Key)
		case "rename":
			if op.NewKey == "" {
				return nil, fmt.Errorf("patch: rename op for key %q has empty new_key", op.Key)
			}
			val, ok := patched.Data[op.Key]
			if !ok {
				return nil, fmt.Errorf("patch: rename source key %q not found", op.Key)
			}
			delete(patched.Data, op.Key)
			patched.Data[op.NewKey] = val
		default:
			return nil, fmt.Errorf("patch: unknown op %q", op.Op)
		}
	}

	// Re-derive the ID from the new data so it stays consistent.
	patched.ID = generateID(patched.Data)
	return patched, nil
}
