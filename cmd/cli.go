package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/d3sw/floop"
)

// CLI is the command line interface to floop
type CLI struct {
	Version    string
	ConfigFile string

	isHelp    bool
	isVersion bool
	debug     bool

	Exec []string               // child process command and args
	Meta map[string]interface{} // context data from command line
}

// Command returns the child process command from the cli args
func (cli *CLI) Command() string {
	if len(cli.Exec) > 0 {
		return cli.Exec[0]
	}
	return ""
}

// Args returns args of the child process from the cli args
func (cli *CLI) Args() []string {
	if len(cli.Exec) > 1 {
		return cli.Exec[1:]
	}
	return []string{}
}

// NewCLI instantiates a new CLI with the command line args.
func NewCLI(version string, args []string) (cli *CLI, err error) {
	cli = &CLI{
		Version:    version,
		ConfigFile: os.Getenv("FLOOP_CONFIG"),
		Meta:       make(map[string]interface{}),
	}
	if cli.ConfigFile == "" {
		cli.ConfigFile = "config.yml"
	}

	for i := 0; i < len(args); i++ {

		switch args[i] {
		case "-exec":
			// Must be the last argument.  Everything after is considered part of the child process
			if i+1 >= len(args) {
				err = errors.New("-exec requires options")
			} else {
				cli.Exec = args[i+1:]
			}
			return

		case "-c":
			i++
			cli.ConfigFile = args[i]
		case "-debug":
			cli.debug = true
		case "-h", "-help", "--help", "--h":
			cli.isHelp = true
		case "-version", "--version":
			cli.isVersion = true
		default:
			// Parse as key=value metadata to be passed in
			if strings.Contains(args[i], "=") {
				if arr := strings.Split(strings.TrimSpace(args[i]), "="); len(arr) > 1 {
					cli.Meta[arr[0]] = arr[1]
				}
			}
		}

	}

	return
}

// Usage prints CLI usage
func (cli *CLI) Usage() {
	fmt.Printf(`
Usage: floop [-c <config_file>] [key=value ...] -exec <command> [args]

floop is a tool to add lifecycle event handlers to any arbitrary process

`)
}

// Run runs floop based on the cli args
func (cli *CLI) Run() (int, error) {
	if cli.isHelp {
		cli.Usage()
		return 0, nil
	} else if cli.isVersion {
		fmt.Println(cli.Version)
		return 0, nil
	} else if cli.debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	}

	return cli.run()
}

func (cli *CLI) run() (int, error) {
	exitCode := -1
	conf, err := floop.LoadConfig(cli.ConfigFile)
	if err != nil {
		return exitCode, err
	}

	if !conf.HasMeta(cli.Meta) {
		return exitCode, fmt.Errorf("required metadata: %v", conf.Meta)
	}

	lifeCycle, err := floop.NewLifecycle(conf)
	if err != nil {
		return exitCode, err
	}

	// Override cli command and args
	cmd := cli.Command()
	if cmd != "" {
		conf.Command = cmd
		conf.Args = cli.Args()
	}

	input := newInput(conf.Command, conf.Args, lifeCycle, conf.Quiet)
	lci, err := floop.NewLifecycledChild(input, lifeCycle)
	if err != nil {
		return exitCode, err
	}
	if err = lci.Start(cli.Meta); err != nil {
		return exitCode, err
	}

	exitCode = lci.Wait()
	return exitCode, nil
}
