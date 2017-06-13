package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/d3sw/floop"
)

// -exec must be the last argument
func parseCLI(args []string) (string, *floop.Context, error) {
	ctx := &floop.Context{
		Meta: make(map[string]interface{}),
	}
	cfgfile := "config.yml"

	for i := 0; i < len(args); i++ {

		if args[i] == "-exec" {
			if i+1 >= len(args) {
				return cfgfile, ctx, errors.New("-exec requires options")
			}

			ctx.Command = args[i+1]
			if i+2 < len(args) {
				ctx.Args = args[i+2:]
			}

			break

		} else if args[i] == "-c" {
			i++
			cfgfile = args[i]
			continue
		}

		if strings.Contains(args[i], "=") {
			if arr := strings.Split(strings.TrimSpace(args[i]), "="); len(arr) > 1 {
				ctx.Meta[arr[0]] = arr[1]
			}
		}

	}

	return cfgfile, ctx, nil
}

func usage() {
	fmt.Printf(`floop [-c <config_file>] [key=value ...] -exec [command] [args]

floop is a tool to add lifecycle event handlers to any arbitrary process

`)
}
