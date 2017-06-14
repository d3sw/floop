package floop

import (
	"fmt"
	"log"

	"github.com/d3sw/floop/handlers"
	"github.com/d3sw/floop/types"
)

// Handler implements and event handler interface.
type Handler interface {
	Handle(event *types.Event) (map[string]interface{}, error)
}

type phaseHandler struct {
	conf *HandlerConfig
	Handler
}

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

// Lifecycle implements a Lifecycle that calls multiple lifecycles for an event.
type Lifecycle struct {
	ctx      *types.Context
	handlers map[types.EventType][]*phaseHandler
}

// NewLifecycle instantiates an instance of Lifecycle
func NewLifecycle(conf *Config) (*Lifecycle, error) {
	lc := &Lifecycle{handlers: make(map[types.EventType][]*phaseHandler)}
	if conf == nil {
		return lc, nil
	}
	err := lc.loadHandlers(conf)
	return lc, err
}

func (lc *Lifecycle) loadHandlers(conf *Config) error {
	for eventType, configs := range conf.Handlers {
		for _, config := range configs {

			var handler Handler

			switch config.Type {
			case "http":
				cfg := &handlers.EndpointConfig{
					URI:     config.Config["uri"].(string),
					Method:  config.Config["method"].(string),
					Headers: make(map[string]string),
				}
				if _, ok := config.Config["body"]; ok {
					cfg.Body = config.Config["body"].(string)
				}

				if hdrs, ok := config.Config["headers"]; ok {
					hm, ok := hdrs.(map[interface{}]interface{})
					if !ok {
						return fmt.Errorf("invalid header data type %#v", config.Config["headers"])
					}
					for k, v := range hm {
						key := k.(string)
						value := v.(string)
						cfg.Headers[key] = value
					}
				}

				handler = handlers.NewHTTPClientHandler(cfg)
			case "echo":
				handler = &handlers.EchoHandler{}
			default:
				return fmt.Errorf("handler not supported: %s", config.Type)
			}

			lc.register(eventType, handler, config)
			log.Printf("[INFO] Registered handler: phase=%s handler=%s", eventType, config.Type)
		}
	}
	return nil
}

// Register registers a new Handler by an arbitrary name.
func (lc *Lifecycle) register(eventType types.EventType, l Handler, conf *HandlerConfig) {
	arr, ok := lc.handlers[eventType]
	if !ok {
		lc.handlers[eventType] = []*phaseHandler{&phaseHandler{Handler: l, conf: conf}}
		return
	}

	lc.handlers[eventType] = append(arr, &phaseHandler{Handler: l, conf: conf})
}

// Begin is called right before a process is launched.  The context is internally stored and may be
// updated by subsequent phases from callback responses.
func (lc *Lifecycle) Begin(ctx *types.Context) error {
	lc.ctx = ctx

	handlers, ok := lc.handlers[types.EventTypeBegin]
	if !ok || handlers == nil || len(handlers) == 0 {
		return nil
	}

	event := &types.Event{Type: types.EventTypeBegin, Meta: ctx.Meta}
	for _, v := range handlers {
		meta, err := v.Handle(event)
		if err != nil {
			if v.conf.IgnoreErrors {
				log.Printf("[ERROR] phase=%s handler=%s %v", event.Type, v.conf.Type, err)
				continue
			}
			return err
		}

		lc.applyContext(meta, v.conf)
	}

	return nil
}

// Progress echos the progress payload
func (lc *Lifecycle) Progress(line []byte) {
	//fmt.Printf("[Progress] %s", line)
	handlers, ok := lc.handlers[types.EventTypeProgress]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	//event := &types.Event{Type: types.EventTypeProgress, Meta: lc.ctx.Meta, Data: line}
	for _, v := range handlers {
		event := &types.Event{Type: types.EventTypeProgress, Meta: lc.ctx.Meta}
		if !v.applyTransform(string(line), event) {
			event.Data = line
		}

		if _, err := v.Handle(event); err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", event.Type, v.conf.Type, err)
		}
	}
}

// Failed is called if the process exits with a non-zero exit status.
func (lc *Lifecycle) Failed(exitCode int) {

	handlers, ok := lc.handlers[types.EventTypeFailed]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	event := &types.Event{Type: types.EventTypeFailed, Meta: lc.ctx.Meta, Data: exitCode}
	for _, v := range handlers {
		if _, err := v.Handle(event); err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", event.Type, v.conf.Type, err)
		}
	}
}

// Completed is called when a process exits with a zero exit code.
func (lc *Lifecycle) Completed() {
	handlers, ok := lc.handlers[types.EventTypeCompleted]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}
	//fmt.Printf("[End] %d\n", exitCode)
	event := &types.Event{Type: types.EventTypeCompleted, Meta: lc.ctx.Meta}
	for _, v := range handlers {
		if _, err := v.Handle(event); err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", event.Type, v.conf.Type, err)
		}
	}
}

func (lc *Lifecycle) applyContext(meta map[string]interface{}, conf *HandlerConfig) {
	if conf.Context == nil || len(conf.Context) == 0 {
		return
	}

	if meta != nil {
		for _, v := range conf.Context {
			if val, ok := meta[v]; ok {
				lc.ctx.Meta[v] = val
			}
		}
	}
}
