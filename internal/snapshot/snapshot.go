package snapshot

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// Snapshot represents a captured state of infrastructure at a point in time.
type Snapshot struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label,omitempty"`
	Entries   map[string]string `json:"entries"`
}

// New creates a new Snapshot with the given label and key-value entries.
func New(label string, entries map[string]string) *Snapshot {
	now := time.Now().UTC()
	return &Snapshot{
		ID:        generateID(now, entries),
		Timestamp: now,
		Label:     label,
		Entries:   entries,
	}
}

// Marshal serializes the snapshot to JSON bytes.
func (s *Snapshot) Marshal() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// Unmarshal deserializes a snapshot from JSON bytes.
func Unmarshal(data []byte) (*Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: failed to unmarshal: %w", err)
	}
	return &s, nil
}

// generateID creates a deterministic SHA-256-based ID from timestamp and entries.
func generateID(t time.Time, entries map[string]string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d", t.UnixNano())
	for k, v := range entries {
		fmt.Fprintf(h, "%s=%s", k, v)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:12]
}
