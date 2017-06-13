package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/d3sw/floop"
	"github.com/d3sw/floop/child"
	"github.com/d3sw/floop/handlers"
	"github.com/d3sw/floop/lifecycle"
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

func newInput(args []string, lc *lifecycle.Lifecycle, noDisplay bool) *child.NewInput {
	stdout := floop.NewBufferedWriter(lc.Progress)
	stderr := floop.NewBufferedWriter(lc.Progress)
	input := commandLessInput(stdout, stderr, noDisplay)

	input.Command = args[0]
	input.Args = args[1:]
	return input
}

func loadHandlers(lifeCycle *lifecycle.Lifecycle, conf *floop.Config) error {
	for eventType, configs := range conf.Handlers {
		for _, config := range configs {

			var handler lifecycle.Handler

			switch config.Type {
			case "http":
				cfg := &handlers.EndpointConfig{
					URI:    config.Config["uri"].(string),
					Method: config.Config["method"].(string),
					Body:   config.Config["body"].(string),
				}
				handler = handlers.NewHTTPClientHandler(cfg)

				//log.Printf("[DEBUG] Registered handler: event=%s handler=%s", eventType, config.Type)
			case "ffmpeg":
				handler = &handlers.FFMPEGHandler{}

			default:
				return fmt.Errorf("handler not supported: %s", config.Type)
			}

			lifeCycle.Register(eventType, handler)
		}
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	//flag.Parse()

	conf, err := floop.LoadConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Conf %v", conf.Handlers.HTTP.EndpointConfig)

	lifeCycle := lifecycle.New()
	err = loadHandlers(lifeCycle, conf)
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]
	input := newInput(args, lifeCycle, conf.Quiet)

	lci, err := floop.NewLifecycledChild(input, lifeCycle)
	if err != nil {
		log.Fatal(err)
	}

	if err = lci.Start(nil); err != nil {
		log.Fatal(err)
	}

	exitCode := lci.Wait()
	// Exit with child process exit status
	os.Exit(exitCode)
}
