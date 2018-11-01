package floop

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/d3sw/floop/child"
	"github.com/d3sw/floop/types"
)

// Floop is the core interface that manages the process lifecycle and handlers
type Floop struct {
	lifecycle *Lifecycle

	bufOut *BufferedWriter // writer to manage progress
	bufErr *BufferedWriter // writer to manage progress

	procInput *child.NewInput
	proc      *child.Child
}

// New instantiates a new instance of floop.
func New(conf *Config, input *child.NewInput) (*Floop, error) {
	lifecycle, err := NewLifecycle(conf)
	if err != nil {
		return nil, err
	}

	var errCallbackWriter func([]byte)
	if conf.ReadFromStderr {
		errCallbackWriter = lifecycle.Progress
	}
	flp := &Floop{
		lifecycle: lifecycle,
		bufOut:    NewBufferedWriter(lifecycle.Progress, true),
		bufErr:    NewBufferedWriter(errCallbackWriter, true),
	}

	input.Command = conf.Command
	input.Args = conf.Args

	input.Stdin = os.Stdin
	if conf.Quiet {
		input.Stdout = flp.bufOut
		input.Stderr = flp.bufErr
	} else {
		input.Stdout = io.MultiWriter(flp.bufOut, os.Stdout)
		input.Stderr = io.MultiWriter(flp.bufErr, os.Stderr)
	}

	flp.procInput = input
	flp.proc, err = child.New(flp.procInput)
	go flp.listenSignals(flp.proc.Signal)

	return flp, err
}

func (floop *Floop) listenSignals(sigProcesser func(os.Signal) error) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)
	sig := <-signalChannel
	log.Printf("[INFO] (floop) got \"%s\" signal\n", sig)
	if err := sigProcesser(sig); err != nil {
		log.Printf("[ERR] (floop) calling signal [%+v] to child: %s\n", sig, err.Error())
	}
}

// Start calls the begin phase of the lifecycle and starts the child process
func (floop *Floop) Start(meta map[string]interface{}) error {
	ctx := &types.Context{
		Command: floop.procInput.Command,
		Args:    floop.procInput.Args,
		Meta:    meta,
	}

	if err := floop.lifecycle.Begin(ctx); err != nil {
		return err
	}
	return floop.proc.Start()
}

// Wait waits for the child process to exit and calls the end phase of the lifecycle
func (floop *Floop) Wait() int {
	code := <-floop.proc.ExitCh()

	result := &types.ChildResult{
		Code:   code,
		Stdout: bytes.TrimRight(floop.bufOut.Bytes(), "\n"),
		Stderr: bytes.TrimRight(floop.bufErr.Bytes(), "\n"),
	}

	if code != 0 {
		state := floop.proc.State().String()
		if strings.Contains(state, os.Interrupt.String()) || strings.Contains(state, os.Kill.String()) {
			floop.lifecycle.Canceled(result)
		} else {
			floop.lifecycle.Failed(result)
		}
	} else {
		floop.lifecycle.Completed(result)
	}

	floop.lifecycle.Close()

	return code
}
