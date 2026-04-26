package alias_test

import (
	"testing"

	"github.com/yourorg/trackpoint/internal/alias"
)

func TestResolve_KnownAlias_ReturnsMappedID(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	_ = st.Set("prod", "snapshot-abc")
	r := alias.NewResolver(st)

	id, err := r.Resolve("prod")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if id != "snapshot-abc" {
		t.Errorf("got %q, want %q", id, "snapshot-abc")
	}
}

func TestResolve_UnknownToken_ReturnsSelf(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	r := alias.NewResolver(st)

	token := "raw-snapshot-id-xyz"
	id, err := r.Resolve(token)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if id != token {
		t.Errorf("got %q, want %q", id, token)
	}
}

func TestResolve_EmptyToken_ReturnsError(t *testing.T) {
	st := alias.NewStore(tempDir(t))
	r := alias.NewResolver(st)

	_, err := r.Resolve("")
	if err == nil {
		t.Error("expected error for empty token")
	}
}
