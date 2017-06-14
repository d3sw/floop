package floop

import "github.com/d3sw/floop/types"

// Handler represents the event handler interface.
type Handler interface {
	Handle(event *types.Event) (map[string]interface{}, error)
	Init(conf *types.HandlerConfig) error
}

// phaseHandler is the internal handler wrapping the config and handler interfaces
type phaseHandler struct {
	conf *types.HandlerConfig
	Handler
}

// applyTransform applies a transform to the input using the event and configuration returning if
// a transform was performed
func (handler *phaseHandler) applyTransform(input string, out *types.Event) bool {
	if len(handler.conf.Transform) == 0 {
		return false
	}

	var transformed bool

	tf := handler.conf.Transform
	switch tf[0] {
	case "kv":
		kvs := transformKeyValuePairs(input, tf[1], tf[2])
		if len(kvs) > 0 {
			out.Data = kvs
			transformed = true
		}
	case "line":
		lines := transformLines(input, tf[1])
		if len(lines) > 0 {
			out.Data = lines
			transformed = true
		}
	}

	return transformed
}
