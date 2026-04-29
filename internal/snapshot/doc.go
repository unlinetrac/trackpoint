// Package snapshot provides types and utilities for capturing, storing,
// and manipulating infrastructure state snapshots.
//
// A Snapshot is an immutable record of key/value state data taken at a
// specific point in time. Each snapshot carries a deterministic content-
// addressed ID derived from its data, a creation timestamp, and an
// optional map of user-defined tags.
//
// # Core Operations
//
//   - New: construct a snapshot from a key/value map
//   - Clone: deep-copy a snapshot, producing a new ID and timestamp
//   - CloneWithOverrides: clone and merge additional key/value pairs
//   - Marshal / Unmarshal: JSON serialization round-trip
//
// # Storage
//
// The Store type persists snapshots to a local directory as JSON files.
// Each file is named by the snapshot ID. Store supports Save, Load, and
// List operations.
//
// Snapshots are intended to be used alongside the diff, filter, history,
// and report packages for change tracking across deploys.
package snapshot
