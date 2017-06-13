package lifecycle

import "log"

// Context is the context passed as part of Lifecycle events
type Context struct {
	Command string
	Args    []string
	Meta    map[string]string
}

// Lifecycle implements a Lifecycle that calls multiple lifecycles for an event.
type Lifecycle struct {
	ctx      *Context
	handlers map[EventType][]Handler
}

// New instantiates an instance of Lifecycle
func New() *Lifecycle {
	return &Lifecycle{handlers: make(map[EventType][]Handler)}
}

// Register registers a new Handler by an arbitrary name.
func (lc *Lifecycle) Register(eventType EventType, l Handler) {
	arr, ok := lc.handlers[eventType]
	if !ok {
		lc.handlers[eventType] = []Handler{l}
		return
	}

	lc.handlers[eventType] = append(arr, l)
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
		if err := v.Handle(event); err != nil {
			log.Println("[ERROR]", event.Type, err)
		}
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
		if err := v.Handle(event); err != nil {
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
		if err := v.Handle(event); err != nil {
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
		if err := v.Handle(event); err != nil {
			log.Println("[ERROR]", event.Type, err)
		}
	}
}
