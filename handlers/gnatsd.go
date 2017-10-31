package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/d3sw/floop/types"
	"github.com/nats-io/go-nats"
)

// GnatsdHandler is handler to publish lifecycle events to gnatsd
type GnatsdHandler struct {
	conf *types.HandlerConfig
	// Nats connection
	conn *nats.Conn
}

// Init initializes the connection to the gnatsd cluster
func (lc *GnatsdHandler) Init(conf *types.HandlerConfig) error {
	lc.conf = conf
	topic, ok := lc.conf.Options.GetString("topic")
	if !ok || topic == "" {
		return fmt.Errorf("topic required or invalid topic: %v", topic)
	}

	opts := nats.Options{
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
		Url:            conf.URI,
		AsyncErrorCB: func(conn *nats.Conn, sub *nats.Subscription, err error) {
			log.Println("[ERROR]", err)
		},
	}

	var err error
	lc.conn, err = opts.Connect()
	return err
}

// Handle publishes to gnatsd.  The config is the normalized
// config built using data from the child process.  This may be different from the one
// used in Init
func (lc *GnatsdHandler) Handle(event *types.Event, conf *types.HandlerConfig) (map[string]interface{}, error) {
	// Get topic from config
	topic, ok := conf.Options.GetString("topic")
	if !ok || topic == "" {
		return nil, fmt.Errorf("topic not specified")
	}

	fmt.Printf("[gnatsd] phase=%s topic=%s %+v\n", event.Type, topic, event.Data)

	// Flushes the connection ensuring the last event gets published before the app terminates
	defer lc.conn.Close()

	// Publish the body as bytes
	err := lc.conn.Publish(topic, []byte(conf.Body))

	// if err == nil {
	// 	if err = lc.econn.Flush(); err != nil {
	// 		err = lc.econn.LastError()
	// 	}
	// }

	return nil, err
}
