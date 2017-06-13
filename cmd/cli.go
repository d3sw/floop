package main

import (
	"errors"
	"strings"

	"github.com/d3sw/floop"
)

func parseCLI(args []string) (*floop.Context, error) {
	ctx := &floop.Context{
		Meta: make(map[string]interface{}),
	}

	for i, arg := range args {

		if arg == "-exec" {
			if i+1 >= len(args) {
				return ctx, errors.New("-exec requires options")
			}

			ctx.Command = args[i+1]
			if i+2 < len(args) {
				ctx.Args = args[i+2:]
			}

			break

		}

		if strings.Contains(arg, "=") {
			if arr := strings.Split(strings.TrimSpace(arg), "="); len(arr) > 1 {
				ctx.Meta[arr[0]] = arr[1]
			}
		}

	}

	return ctx, nil
}
