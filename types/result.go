package types

// ChildResult is the result of process
type ChildResult struct {
	Code   int // exit code
	Stdout []byte
	Stderr []byte
}
