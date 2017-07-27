package types

// Context is the context passed as part of Lifecycle events
type Context struct {
	Command string
	Args    []string
	Meta    map[string]interface{}
}
