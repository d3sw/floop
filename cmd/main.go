package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/d3sw/floop"
	"github.com/d3sw/floop/child"
	"github.com/d3sw/floop/handlers"
)

//var (
// onBegin    = flag.String("begin", "", "")
// onProgress = flag.String("progress", "", "")
// onEnd      = flag.String("end", "", "")
//noDisplay = flag.Bool("no-display", false, "Do not write to stdout or stderr")
//)

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

func loadHandlers(lifeCycle *floop.Lifecycle, conf *floop.Config) error {
	for eventType, configs := range conf.Handlers {
		for _, config := range configs {

			var handler floop.Handler

			switch config.Type {
			case "http":
				cfg := &handlers.EndpointConfig{
					URI:    config.Config["uri"].(string),
					Method: config.Config["method"].(string),
				}
				if _, ok := config.Config["body"]; ok {
					cfg.Body = config.Config["body"].(string)
				}
				handler = handlers.NewHTTPClientHandler(cfg)

				//log.Printf("[DEBUG] Registered handler: event=%s handler=%s", eventType, config.Type)
			case "ffmpeg":
				handler = &handlers.FFMPEGHandler{}

			default:
				return fmt.Errorf("handler not supported: %s", config.Type)
			}

			lifeCycle.Register(eventType, handler, config)
		}
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	//flag.Parse()

	args := os.Args[1:]

	ctx, err := parseCLI(args)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := floop.LoadConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	lifeCycle := floop.NewLifecycle()
	err = loadHandlers(lifeCycle, conf)
	if err != nil {
		log.Fatal(err)
	}

	input := newInput(ctx.Command, ctx.Args, lifeCycle, conf.Quiet)

	lci, err := floop.NewLifecycledChild(input, lifeCycle)
	if err != nil {
		log.Fatal(err)
	}

	if err = lci.Start(ctx.Meta); err != nil {
		log.Fatal(err)
	}

	exitCode := lci.Wait()
	// Exit with child process exit status
	os.Exit(exitCode)
}
