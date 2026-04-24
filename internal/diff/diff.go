package diff

// ChangeType describes the nature of a single key-level change.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single key-level difference between two snapshots.
type Change struct {
	Key  string
	Type ChangeType
	From string
	To   string
}

// Result holds the full output of comparing two snapshots.
type Result struct {
	FromID  string
	ToID    string
	Changes []Change
}

// HasChanges returns true when the result contains at least one change.
func (r Result) HasChanges() bool {
	return len(r.Changes) > 0
}

type snapshotIface interface {
	GetID() string
	GetData() map[string]string
}

// Compare produces a Result describing all differences between from and to.
func Compare(from, to snapshotIface) Result {
	result := Result{
		FromID: from.GetID(),
		ToID:   to.GetID(),
	}

	fromData := from.GetData()
	toData := to.GetData()

	for k, fv := range fromData {
		if tv, ok := toData[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Type: Removed, From: fv})
		} else if fv != tv {
			result.Changes = append(result.Changes, Change{Key: k, Type: Modified, From: fv, To: tv})
		}
	}

	for k, tv := range toData {
		if _, ok := fromData[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Type: Added, To: tv})
		}
	}

	return result
}
