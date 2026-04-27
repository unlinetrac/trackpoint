package snapshot

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Snapshot holds a point-in-time capture of key/value infrastructure state.
type Snapshot struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Data      map[string]string `json:"data"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// New creates a new Snapshot from the provided data map.
func New(data map[string]string) *Snapshot {
	now := time.Now().UTC()
	return &Snapshot{
		ID:        generateID(data, now),
		CreatedAt: now,
		Data:      data,
		Tags:      make(map[string]string),
	}
}

// Marshal serialises a Snapshot to JSON bytes.
func (s *Snapshot) Marshal() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Unmarshal deserialises JSON bytes into a Snapshot.
func Unmarshal(b []byte) (*Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// generateID produces a deterministic SHA-256-based ID from the snapshot data
// and creation timestamp.
func generateID(data map[string]string, t time.Time) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	fmt.Fprintf(h, "%d", t.UnixNano())
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, data[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}
