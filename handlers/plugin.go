package handlers

import (
	"fmt"

	"github.com/d3sw/floop/types"
)

type Handler interface {
	Handle(topic string, data []byte) (map[string]interface{}, error)
	Init(map[string]string) error
	CloseConnection() error
}

// PluginHandler is handler to publish lifecycle events
type PluginHandler struct {
	h Handler
}

func NewPluginHandler(handler Handler) *PluginHandler {
	return &PluginHandler{
		h: handler,
	}
}

// Init initializes the connection
func (lc *PluginHandler) Init(conf *types.HandlerConfig) error {
	c := map[string]string{}
	if conf.URI == "" {
		return fmt.Errorf("uri not specified")
	}
	c["uri"] = conf.URI
	return lc.h.Init(c)
}

// Handle publishes messages
func (lc *PluginHandler) Handle(event *types.Event, conf *types.HandlerConfig) (map[string]interface{}, error) {
	// Get topic from config
	topic, ok := conf.Options.GetString("topic")
	if !ok || topic == "" {
		return nil, fmt.Errorf("topic not specified")
	}

	fmt.Printf("[nats-stream] phase=%s topic=%s %+v\n", event.Type, topic, event.Data)

	// Publish the body as bytes
	return lc.h.Handle(topic, []byte(conf.Body))
}

// CloseConnection closes connection
func (lc *PluginHandler) CloseConnection() error {
	return lc.h.CloseConnection()
}
