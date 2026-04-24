package report

import (
	"fmt"
	"io"
	"time"

	"github.com/trackpoint/internal/diff"
	"github.com/trackpoint/internal/snapshot"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a rendered comparison between two snapshots.
type Report struct {
	From      string         `json:"from"`
	To        string         `json:"to"`
	Timestamp time.Time      `json:"timestamp"`
	Result    diff.Result    `json:"result"`
}

// New creates a Report by comparing two snapshots.
func New(from, to *snapshot.Snapshot) (*Report, error) {
	if from == nil {
		return nil, fmt.Errorf("report: from snapshot is nil")
	}
	if to == nil {
		return nil, fmt.Errorf("report: to snapshot is nil")
	}

	result := diff.Compare(from, to)

	return &Report{
		From:      from.ID,
		To:        to.ID,
		Timestamp: time.Now().UTC(),
		Result:    result,
	}, nil
}

// Write renders the report to the given writer in the specified format.
func (r *Report) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, r)
	case FormatText:
		return writeText(w, r)
	default:
		return fmt.Errorf("report: unsupported format %q", format)
	}
}

// HasChanges returns true if the report contains any diff entries.
func (r *Report) HasChanges() bool {
	return len(r.Result.Changes) > 0
}
