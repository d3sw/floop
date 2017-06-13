package floop

import (
	"log"
)

type phaseHandler struct {
	conf *HandlerConfig
	Handler
}

// Context is the context passed as part of Lifecycle events
type Context struct {
	Command string
	Args    []string
	Meta    map[string]interface{}
}

// Lifecycle implements a Lifecycle that calls multiple lifecycles for an event.
type Lifecycle struct {
	ctx      *Context
	handlers map[EventType][]*phaseHandler
}

// NewLifecycle instantiates an instance of Lifecycle
func NewLifecycle() *Lifecycle {
	return &Lifecycle{handlers: make(map[EventType][]*phaseHandler)}
}

// Register registers a new Handler by an arbitrary name.
func (lc *Lifecycle) Register(eventType EventType, l Handler, conf *HandlerConfig) {
	arr, ok := lc.handlers[eventType]
	if !ok {
		lc.handlers[eventType] = []*phaseHandler{&phaseHandler{Handler: l, conf: conf}}
		return
	}

	lc.handlers[eventType] = append(arr, &phaseHandler{Handler: l, conf: conf})
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

// Begin echos back input data before process starts
func (lc *Lifecycle) Begin(ctx *Context) {
	lc.ctx = ctx

	handlers, ok := lc.handlers[EventTypeBegin]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}
	//fmt.Printf("[Begin] %+v\n", lc.ctx)
	event := &Event{Type: EventTypeBegin, Meta: ctx.Meta}

	for _, v := range handlers {

		meta, err := v.Handle(event)
		if err != nil {
			log.Println("[ERROR]", event.Type, err)
			continue
		}

		lc.applyContext(meta, v.conf)

	}
}

// Progress echos the progress payload
func (lc *Lifecycle) Progress(line []byte) {
	//fmt.Printf("[Progress] %s", line)
	handlers, ok := lc.handlers[EventTypeProgress]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	event := &Event{Type: EventTypeProgress, Meta: lc.ctx.Meta, Data: line}
	for _, v := range handlers {
		if _, err := v.Handle(event); err != nil {
			log.Println("[ERROR]", event.Type, err)
		}
	}
}

// Failed is called if the process exits with a non-zero exit status.
func (lc *Lifecycle) Failed(exitCode int) {

	handlers, ok := lc.handlers[EventTypeFailed]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}
	//fmt.Printf("[End] %d\n", exitCode)
	event := &Event{Type: EventTypeFailed, Meta: lc.ctx.Meta, Data: exitCode}
	for _, v := range handlers {
		if _, err := v.Handle(event); err != nil {
			log.Println("[ERROR]", event.Type, err)
		}
	}
}

// Completed is called when a process exits with a zero exit code.
func (lc *Lifecycle) Completed() {
	handlers, ok := lc.handlers[EventTypeCompleted]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}
	//fmt.Printf("[End] %d\n", exitCode)
	event := &Event{Type: EventTypeCompleted, Meta: lc.ctx.Meta}
	for _, v := range handlers {
		if _, err := v.Handle(event); err != nil {
			log.Println("[ERROR]", event.Type, err)
		}
	}
}
