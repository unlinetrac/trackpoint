package snapshot

import (
	"encoding/json"
	"fmt"
)

// Clone returns a deep copy of the given Snapshot with a new ID and timestamp.
// The cloned snapshot inherits all data and tags from the original.
func Clone(s *Snapshot) (*Snapshot, error) {
	if s == nil {
		return nil, fmt.Errorf("snapshot: cannot clone nil snapshot")
	}

	// Deep copy data map
	dataCopy := make(map[string]string, len(s.Data))
	for k, v := range s.Data {
		dataCopy[k] = v
	}

	// Deep copy tags map
	tagsCopy := make(map[string]string, len(s.Tags))
	for k, v := range s.Tags {
		tagsCopy[k] = v
	}

	cloned := New(dataCopy)
	cloned.Tags = tagsCopy
	return cloned, nil
}

// CloneWithOverrides returns a deep copy of the snapshot with the provided
// data keys merged on top. Existing keys not in overrides are preserved.
func CloneWithOverrides(s *Snapshot, overrides map[string]string) (*Snapshot, error) {
	cloned, err := Clone(s)
	if err != nil {
		return nil, err
	}
	for k, v := range overrides {
		cloned.Data[k] = v
	}
	// Recompute ID based on new data
	raw, err := json.Marshal(cloned.Data)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to marshal overridden data: %w", err)
	}
	cloned.ID = generateID(raw)
	return cloned, nil
}
