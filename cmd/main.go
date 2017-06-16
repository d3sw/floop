package main

import (
	"fmt"
	"os"
)

var version = "unknown"

func main() {

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
