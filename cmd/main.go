package main

import (
	"fmt"
	"io"
	"os"

	"github.com/d3sw/floop"
	"github.com/d3sw/floop/child"
)

var version = "unknown"

func commandLessInput(sout, serr io.Writer, noDisplay bool) *child.NewInput {
	var (
		stdout io.Writer
		stderr io.Writer
	)

	if noDisplay {
		stdout = sout
		stderr = serr
	} else {
		stdout = io.MultiWriter(os.Stdout, sout)
		stderr = io.MultiWriter(os.Stderr, serr)
	}

	return &child.NewInput{
		Stdin:  os.Stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func newInput(cmd string, args []string, lc *floop.Lifecycle, noDisplay bool) *child.NewInput {
	stdout := floop.NewBufferedWriter(lc.Progress)
	stderr := floop.NewBufferedWriter(lc.Progress)
	input := commandLessInput(stdout, stderr, noDisplay)

	input.Command = cmd
	input.Args = args
	return input
}

func main() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	//flag.Parse()

	args := os.Args[1:]
	cli, err := NewCLI(version, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Exit with child process exit status
	os.Exit(exitCode)
}
