package alias

import (
	"errors"
	"fmt"
)

// Resolver resolves a token that may be either a raw snapshot ID or an alias
// name into a concrete snapshot ID.
type Resolver struct {
	store *Store
}

// NewResolver creates a Resolver backed by the given Store.
func NewResolver(store *Store) *Resolver {
	return &Resolver{store: store}
}

// Resolve returns the snapshot ID for token. If token is a known alias the
// mapped ID is returned; otherwise token itself is returned unchanged,
// allowing callers to pass raw IDs without needing special handling.
func (r *Resolver) Resolve(token string) (string, error) {
	if token == "" {
		return "", errors.New("token must not be empty")
	}
	id, err := r.store.Get(token)
	if err == nil {
		return id, nil
	}
	if errors.Is(err, ErrNotFound) {
		// Treat as a raw snapshot ID.
		return token, nil
	}
	return "", fmt.Errorf("resolve alias %q: %w", token, err)
}
