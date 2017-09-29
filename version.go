package floop

import "fmt"

// VERSION number
const VERSION = "0.1.1"

var (
	branch    string
	commit    string
	buildtime string
)

func printVersion() {
	fmt.Printf("%s commit=%s/%s buildtime=%s\n", VERSION, branch, commit, buildtime)
}
