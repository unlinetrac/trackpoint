package watch

import (
	"time"

	"github.com/user/trackpoint/internal/diff"
	"github.com/user/trackpoint/internal/snapshot"
)

// Config holds configuration for a watch session.
type Config struct {
	Interval  time.Duration
	MaxChecks int // 0 means run indefinitely
}

// Change represents a detected change event during a watch.
type Change struct {
	At     time.Time
	Result diff.Result
}

// Watcher polls a snapshot source at a given interval and emits changes.
type Watcher struct {
	cfg      Config
	store    *snapshot.Store
	previous *snapshot.Snapshot
}

// New creates a Watcher using the provided store and config.
func New(store *snapshot.Store, cfg Config) *Watcher {
	if cfg.Interval <= 0 {
		cfg.Interval = 10 * time.Second
	}
	return &Watcher{cfg: cfg, store: store}
}

// Watch polls for new snapshots and sends Change events on the returned channel.
// The channel is closed when the watch ends or ctx is cancelled.
func (w *Watcher) Watch(snapshotIDs []string, out chan<- Change) error {
	defer close(out)

	ids := snapshotIDs
	checks := 0

	for {
		if w.cfg.MaxChecks > 0 && checks >= w.cfg.MaxChecks {
			return nil
		}

		for i := 1; i < len(ids); i++ {
			from, err := w.store.Load(ids[i-1])
			if err != nil {
				return err
			}
			to, err := w.store.Load(ids[i])
			if err != nil {
				return err
			}
			result := diff.Compare(from, to)
			if len(result.Changes) > 0 {
				out <- Change{At: time.Now(), Result: result}
			}
		}

		checks++
		if w.cfg.MaxChecks > 0 && checks >= w.cfg.MaxChecks {
			return nil
		}
		time.Sleep(w.cfg.Interval)
	}
}
