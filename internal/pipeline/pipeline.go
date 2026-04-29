// Package pipeline provides a composable, ordered sequence of snapshot
// transformations applied before diffing or reporting.
package pipeline

import (
	"errors"
	"fmt"

	"github.com/user/trackpoint/internal/snapshot"
)

// StageFunc is a function that transforms a snapshot, returning the modified
// copy or an error that halts the pipeline.
type StageFunc func(snap *snapshot.Snapshot) (*snapshot.Snapshot, error)

// Stage is a named transformation step.
type Stage struct {
	Name string
	Fn   StageFunc
}

// Pipeline is an ordered list of stages applied sequentially to a snapshot.
type Pipeline struct {
	stages []Stage
}

// New returns an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a stage to the pipeline.
func (p *Pipeline) Add(name string, fn StageFunc) *Pipeline {
	if name == "" {
		panic("pipeline: stage name must not be empty")
	}
	if fn == nil {
		panic("pipeline: stage function must not be nil")
	}
	p.stages = append(p.stages, Stage{Name: name, Fn: fn})
	return p
}

// Run executes each stage in order, passing the output of one stage as the
// input to the next. It returns the final snapshot or the first error
// encountered, annotated with the stage name.
func (p *Pipeline) Run(snap *snapshot.Snapshot) (*snapshot.Snapshot, error) {
	if snap == nil {
		return nil, errors.New("pipeline: input snapshot must not be nil")
	}
	current := snap
	for _, s := range p.stages {
		out, err := s.Fn(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %q: %w", s.Name, err)
		}
		if out == nil {
			return nil, fmt.Errorf("pipeline stage %q: returned nil snapshot", s.Name)
		}
		current = out
	}
	return current, nil
}

// Stages returns a copy of the registered stage list.
func (p *Pipeline) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}
