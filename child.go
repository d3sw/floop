package floop

import (
	"github.com/d3sw/floop/child"
	"github.com/d3sw/floop/lifecycle"
)

// LifecycledChild wraps lifecycle events around a child process. It hooks in to allow process
// progress via attaching to stdout and stdin.
type LifecycledChild struct {
	input *child.NewInput
	proc  *child.Child

	lc *lifecycle.Lifecycle
}

// NewLifecycledChild instantiates a new LifecycledChild with an input and lifecycle to use for the
// instannce
func NewLifecycledChild(input *child.NewInput, lifecycle *lifecycle.Lifecycle) (*LifecycledChild, error) {

	chld, err := child.New(input)
	if err == nil {
		return &LifecycledChild{
			lc:    lifecycle,
			input: input,
			proc:  chld,
		}, nil
	}

	return nil, err
}

// Start calls the begin phase of the lifecycle and starts the child process
func (li *LifecycledChild) Start(meta map[string]string) error {
	ctx := &lifecycle.Context{
		Command: li.input.Command,
		Args:    li.input.Args,
		Meta:    meta,
	}

	li.lc.Begin(ctx)
	return li.proc.Start()
}

// Wait waits for the child process to exit and calls the end phase of the lifecycle
func (li *LifecycledChild) Wait() int {
	code := <-li.proc.ExitCh()
	if code != 0 {
		li.lc.Failed(code)
	} else {
		li.lc.Completed()
	}
	return code
}
