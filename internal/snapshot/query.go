package snapshot

import (
	"strings"
	"time"
)

// QueryOptions defines filters for selecting snapshots from a store.
type QueryOptions struct {
	// Label filters snapshots whose label contains this substring (case-insensitive).
	Label string

	// Since filters snapshots created at or after this time.
	Since *time.Time

	// Until filters snapshots created before or at this time.
	Until *time.Time

	// Tags filters snapshots that contain all specified key=value tags.
	Tags map[string]string

	// Limit caps the number of results returned (0 = no limit).
	Limit int
}

// Query returns snapshots from the store that match the given options.
// Results are returned in the same order as the store lists them (sorted by ID).
func Query(store Loader, opts QueryOptions) ([]*Snapshot, error) {
	ids, err := store.List()
	if err != nil {
		return nil, err
	}

	var results []*Snapshot
	for _, id := range ids {
		snap, err := store.Load(id)
		if err != nil {
			return nil, err
		}
		if !matchesQuery(snap, opts) {
			continue
		}
		results = append(results, snap)
		if opts.Limit > 0 && len(results) >= opts.Limit {
			break
		}
	}
	return results, nil
}

// Loader is the subset of Store used by Query.
type Loader interface {
	List() ([]string, error)
	Load(id string) (*Snapshot, error)
}

func matchesQuery(snap *Snapshot, opts QueryOptions) bool {
	if opts.Label != "" {
		if !strings.Contains(strings.ToLower(snap.Label), strings.ToLower(opts.Label)) {
			return false
		}
	}
	if opts.Since != nil && snap.CreatedAt.Before(*opts.Since) {
		return false
	}
	if opts.Until != nil && snap.CreatedAt.After(*opts.Until) {
		return false
	}
	for k, v := range opts.Tags {
		if snap.Tags[k] != v {
			return false
		}
	}
	return true
}
