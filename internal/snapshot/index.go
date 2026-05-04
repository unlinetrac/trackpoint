// Package snapshot provides types and utilities for capturing and managing
// point-in-time infrastructure state.
package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// Index provides an in-memory lookup structure for snapshots by label key/value pairs.
type Index struct {
	// labelIndex maps "key=value" -> list of snapshot IDs
	labelIndex map[string][]string
	// idIndex maps snapshot ID -> snapshot
	idIndex map[string]*Snapshot
}

// NewIndex builds an Index from a slice of snapshots.
func NewIndex(snaps []*Snapshot) *Index {
	idx := &Index{
		labelIndex: make(map[string][]string),
		idIndex:    make(map[string]*Snapshot),
	}
	for _, s := range snaps {
		if s == nil {
			continue
		}
		idx.idIndex[s.ID] = s
		for k, v := range s.Labels {
			key := fmt.Sprintf("%s=%s", k, v)
			idx.labelIndex[key] = append(idx.labelIndex[key], s.ID)
		}
	}
	return idx
}

// FindByLabel returns all snapshots matching the given label key and value.
func (idx *Index) FindByLabel(key, value string) []*Snapshot {
	tok := fmt.Sprintf("%s=%s", key, value)
	ids := idx.labelIndex[tok]
	result := make([]*Snapshot, 0, len(ids))
	for _, id := range ids {
		if s, ok := idx.idIndex[id]; ok {
			result = append(result, s)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

// FindByID returns a snapshot by its exact ID, or nil if not found.
func (idx *Index) FindByID(id string) *Snapshot {
	return idx.idIndex[id]
}

// FindByLabelKey returns all snapshots that have the given label key (any value).
func (idx *Index) FindByLabelKey(key string) []*Snapshot {
	seen := make(map[string]struct{})
	var result []*Snapshot
	prefix := key + "="
	for tok, ids := range idx.labelIndex {
		if !strings.HasPrefix(tok, prefix) {
			continue
		}
		for _, id := range ids {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			if s, ok := idx.idIndex[id]; ok {
				result = append(result, s)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}

// Size returns the number of snapshots in the index.
func (idx *Index) Size() int {
	return len(idx.idIndex)
}
