package floop

import (
	"errors"
	"fmt"
	"log"
	"time"

	"plugin"

	"github.com/d3sw/floop/handlers"
	"github.com/d3sw/floop/resolver"
	"github.com/d3sw/floop/types"
)

var (
	dResolverPort = 8600
	dResolverHost = "consul.service"
)

// HandlerResponse stores a result of Handler execution
type handlerResponse struct {
	meta      map[string]interface{}
	err       error
	eType     types.EventType
	conf      *types.HandlerConfig
	ignoreErr bool
}

// Lifecycle implements a Lifecycle that calls multiple lifecycles for an event.
type Lifecycle struct {
	ctx          *types.Context
	handlers     map[types.EventType][]*phaseHandler
	addrResolver *resolver.Resolver
}

// NewLifecycle instantiates an instance of Lifecycle
func NewLifecycle(conf *Config) (*Lifecycle, error) {
	rHosts := conf.ResolverHosts
	rPort := conf.ResolverPort

	if rHosts == nil {
		rHosts = []string{dResolverHost}
	}

	if rPort == 0 {
		rPort = dResolverPort
	}

	lc := &Lifecycle{
		handlers:     make(map[types.EventType][]*phaseHandler),
		addrResolver: resolver.NewResolver(rPort, rHosts...),
	}
	if conf == nil {
		return lc, nil
	}

	err := lc.loadHandlers(conf)
	return lc, err
}

func (lc *Lifecycle) loadHandlers(conf *Config) error {
	for eventType, configs := range conf.Handlers {
		// Setup handlers for an event type
		for _, config := range configs {

			var handler Handler

			switch config.Type {
			case "http":
				interval := 0
				retries := 0

				if _interval, ok := config.Options["interval"]; ok {
					interval = _interval.(int)
				}

				if _retries, ok := config.Options["retries"]; ok {
					retries = _retries.(int)
				}

				var backoff handlers.Backoff
				if _backoff, ok := config.Options["backoff"]; ok && _backoff == "linear" {
					backoff = handlers.LinearBackoff{Interval: time.Duration(interval) * time.Second}
				} else {
					backoff = handlers.ConstantBackoff{Interval: time.Duration(interval) * time.Second}
				}

				handler = handlers.NewHTTPClientHandler(lc.addrResolver, backoff, retries)
			case "echo":
				handler = &handlers.EchoHandler{}
			case "gnatsd":
				handler = &handlers.GnatsdHandler{}
			case "nats-stream":
				handler = &handlers.NatsStreamdHandler{}
			case "plug-in":
				var err error
				if handler, err = lc.loadPluginHandler(config); err != nil {
					return err
				}
			default:
				return fmt.Errorf("handler not supported: %s", config.Type)
			}

			if err := lc.register(eventType, handler, config); err != nil {
				return err
			}
			log.Printf("[INFO] Registered handler: phase=%s handler=%s", eventType, config.Type)
		}
	}
	return nil
}

func (lc *Lifecycle) loadPluginHandler(conf *types.HandlerConfig) (Handler, error) {
	path, ok := conf.Options.GetString("plugin_path")
	if !ok || path == "" {
		return nil, errors.New("plugin_path required")
	}
	name, ok := conf.Options.GetString("plugin_name")
	if !ok || path == "" {
		return nil, errors.New("plugin_name required")
	}

	plugmod, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	plug, err := plugmod.Lookup(name)
	if err != nil {
		return nil, err
	}
	var handler Handler
	handler, ok = plug.(Handler)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}
	return handler, nil
}

// Register registers a new Handler by an arbitrary name.
func (lc *Lifecycle) register(eventType types.EventType, l Handler, conf *types.HandlerConfig) error {
	if err := l.Init(conf); err != nil {
		return err
	}

	arr, ok := lc.handlers[eventType]
	if !ok {
		lc.handlers[eventType] = []*phaseHandler{&phaseHandler{Handler: l, conf: conf}}
	} else {
		lc.handlers[eventType] = append(arr, &phaseHandler{Handler: l, conf: conf})
	}

	return nil
}

// Begin is called right before a process is launched.  The context is internally stored and may be
// updated by subsequent phases from callback responses.
func (lc *Lifecycle) Begin(ctx *types.Context) error {
	lc.ctx = ctx

	handlers, ok := lc.handlers[types.EventTypeBegin]
	if !ok || handlers == nil || len(handlers) == 0 {
		return nil
	}

	doneCh := make(chan *handlerResponse, len(handlers))

	for _, v := range handlers {
		event := &types.Event{
			Type:      types.EventTypeBegin,
			Meta:      ctx.Meta,
			Timestamp: time.Now().UnixNano(),
		}

		go func(v *phaseHandler) {
			meta, err := v.Handle(event)
			doneCh <- &handlerResponse{meta: meta, err: err, eType: event.Type, conf: v.conf, ignoreErr: v.conf.IgnoreErrors}
		}(v)
	}

	var failedHandlers int
	for i := 0; i < len(handlers); i++ {
		res := <-doneCh
		if res.err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", res.eType, res.conf.Type, res.err)
			if res.ignoreErr {
				continue
			}
			failedHandlers++
		}
		lc.applyContext(res.meta, res.conf)
	}

	if failedHandlers > 0 {
		return fmt.Errorf("%d handlers were finished with error", failedHandlers)
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

	doneCh := make(chan *handlerResponse, len(handlers))

	for _, v := range handlers {
		event := &types.Event{
			Type:      types.EventTypeProgress,
			Meta:      lc.ctx.Meta,
			Data:      line,
			Timestamp: time.Now().UnixNano(),
		}
		go func(v *phaseHandler) {
			_, err := v.Handle(event)
			doneCh <- &handlerResponse{err: err, eType: event.Type, conf: v.conf}
		}(v)
	}

	for i := 0; i < len(handlers); i++ {
		res := <-doneCh
		if res.err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", res.eType, res.conf.Type, res.err)
		}
	}
}

// Failed is called if the process exits with a non-zero exit status. Data from stderr and stdout
// are passed in as args
func (lc *Lifecycle) Failed(result *types.ChildResult) {

	handlers, ok := lc.handlers[types.EventTypeFailed]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	doneCh := make(chan *handlerResponse, len(handlers))

	for _, v := range handlers {

		event := &types.Event{
			Type:      types.EventTypeFailed,
			Meta:      lc.ctx.Meta,
			Data:      result,
			Timestamp: time.Now().UnixNano(),
		}
		go func(v *phaseHandler) {
			_, err := v.Handle(event)
			doneCh <- &handlerResponse{err: err, eType: event.Type, conf: v.conf}
		}(v)
	}

	for i := 0; i < len(handlers); i++ {
		res := <-doneCh
		if res.err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", res.eType, res.conf.Type, res.err)
		}
	}
}

// Canceled is called if the process was interrupted or killed. Data from stderr and stdout
// are passed in as args
func (lc *Lifecycle) Canceled(result *types.ChildResult) {
	log.Println("[DEBUG] In Canceled")
	handlers, ok := lc.handlers[types.EventTypeCanceled]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	for _, v := range handlers {

		event := &types.Event{
			Type:      types.EventTypeCanceled,
			Meta:      lc.ctx.Meta,
			Data:      result,
			Timestamp: time.Now().UnixNano(),
		}

		if _, err := v.Handle(event); err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", event.Type, v.conf.Type, err)
		}

	}
}

// Completed is called when a process completes with a zero exit code. Data from stderr and stdout
// are passed in as args
func (lc *Lifecycle) Completed(result *types.ChildResult) {
	handlers, ok := lc.handlers[types.EventTypeCompleted]
	if !ok || handlers == nil || len(handlers) == 0 {
		return
	}

	doneCh := make(chan *handlerResponse, len(handlers))

	for _, v := range handlers {

		event := &types.Event{
			Type:      types.EventTypeCompleted,
			Meta:      lc.ctx.Meta,
			Data:      result.Stdout,
			Timestamp: time.Now().UnixNano(),
		}

		//if !v.applyTransform(result.Stdout, event) {
		//	event.Data = result.Stdout
		//}

		go func(v *phaseHandler) {
			_, err := v.Handle(event)
			doneCh <- &handlerResponse{err: err, eType: event.Type, conf: v.conf}
		}(v)
	}

	for i := 0; i < len(handlers); i++ {
		res := <-doneCh
		if res.err != nil {
			log.Printf("[ERROR] phase=%s handler=%s %v", res.eType, res.conf.Type, res.err)
		}
	}
}

func (lc *Lifecycle) applyContext(meta map[string]interface{}, conf *types.HandlerConfig) {
	if conf.Context == nil || len(conf.Context) == 0 {
		return
	}

	if meta == nil {
		return
	}

	for _, v := range conf.Context {
		if val, ok := meta[v]; ok {
			lc.ctx.Meta[v] = val
		}
	}

}

// Close - close completed and failed handlers so the last event is published
func (lc *Lifecycle) Close() {
	handlers, ok := lc.handlers[types.EventTypeCompleted]
	if ok && handlers != nil && len(handlers) != 0 {
		for _, v := range handlers {
			v.CloseConnection()
		}
	}

	handlers, ok = lc.handlers[types.EventTypeFailed]
	if ok && handlers != nil && len(handlers) != 0 {
		for _, v := range handlers {
			v.CloseConnection()
		}
	}

}
