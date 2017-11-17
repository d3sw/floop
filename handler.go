package floop

import (
	"github.com/d3sw/floop/types"
	"github.com/persephony/shml"
)

// Handler represents the event handler interface.
type Handler interface {
	// raw event and config after interpolation
	Handle(event *types.Event, conf *types.HandlerConfig) (map[string]interface{}, error)
	Init(conf *types.HandlerConfig) error
}

// phaseHandler is the internal handler wrapping the config and handler interfaces
type phaseHandler struct {
	conf *types.HandlerConfig
	Handler
}

func (handler *phaseHandler) buildConfig(event *types.Event) (*types.HandlerConfig, error) {
	// Clone existing config
	conf := handler.conf.Clone()

	// Interpolate Body
	bodyTmpl := shml.New()
	bodyTmpl.Parse([]byte(conf.Body))
	out, err := bodyTmpl.Execute(event)
	if err != nil {
		return nil, err
	}
	conf.Body = string(out)

	// Interpolate URI
	if handler.conf.URI == "" {
		return conf, nil
	}
	uriTmpl := shml.New()
	uriTmpl.Parse([]byte(handler.conf.URI))
	if out, err = uriTmpl.Execute(event); err != nil {
		return nil, err
	}
	conf.URI = string(out)

	return conf, nil
}

func (handler *phaseHandler) Handle(event *types.Event) (map[string]interface{}, error) {
	// Apply transform to the event data before calling the handler.  It is only applied if the
	// data is a byte slice.
	if len(handler.conf.Transform) > 0 {

		if data, ok := event.Data.([]byte); ok {
			if _, err := Transform(handler.conf.Transform, data, event); err != nil {
				return nil, err
			}
		} else if len(event.Data.(*types.ChildResult).Stderr) > 0 || len(event.Data.(*types.ChildResult).Stdout) > 0 {
			if _, err := TransformResult(handler.conf.Transform, event.Data.(*types.ChildResult), event); err != nil {
				return nil, err
			}
		}

	}
	// Build a normalized config to pass to the handler
	conf, err := handler.buildConfig(event)
	if err != nil {
		return nil, err
	}

	// Call user defined handler
	return handler.Handler.Handle(event, conf)
}
