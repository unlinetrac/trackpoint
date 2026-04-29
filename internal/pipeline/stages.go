package pipeline

import (
	"fmt"
	"strings"

	"github.com/user/trackpoint/internal/snapshot"
)

// RedactKeys returns a StageFunc that replaces the values of the given keys
// with the placeholder string (e.g. "***").
func RedactKeys(keys []string, placeholder string) StageFunc {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	return func(snap *snapshot.Snapshot) (*snapshot.Snapshot, error) {
		data := cloneData(snap.Data)
		for k := range set {
			if _, ok := data[k]; ok {
				data[k] = placeholder
			}
		}
		return snapshot.New(snap.Label, data), nil
	}
}

// FilterKeyPrefix returns a StageFunc that removes any key not starting with
// the given prefix.
func FilterKeyPrefix(prefix string) StageFunc {
	return func(snap *snapshot.Snapshot) (*snapshot.Snapshot, error) {
		data := make(map[string]string)
		for k, v := range snap.Data {
			if strings.HasPrefix(k, prefix) {
				data[k] = v
			}
		}
		return snapshot.New(snap.Label, data), nil
	}
}

// NormalizeValues returns a StageFunc that applies a normalizer function to
// every value in the snapshot.
func NormalizeValues(normalize func(string) string) StageFunc {
	if normalize == nil {
		panic("pipeline: normalize function must not be nil")
	}
	return func(snap *snapshot.Snapshot) (*snapshot.Snapshot, error) {
		data := make(map[string]string, len(snap.Data))
		for k, v := range snap.Data {
			data[k] = normalize(v)
		}
		return snapshot.New(snap.Label, data), nil
	}
}

// RequireKeys returns a StageFunc that returns an error if any of the required
// keys are missing from the snapshot.
func RequireKeys(keys []string) StageFunc {
	return func(snap *snapshot.Snapshot) (*snapshot.Snapshot, error) {
		for _, k := range keys {
			if _, ok := snap.Data[k]; !ok {
				return nil, fmt.Errorf("required key %q is missing", k)
			}
		}
		return snap, nil
	}
}

func cloneData(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
