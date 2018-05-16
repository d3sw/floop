package handlers

import (
	"errors"
	"fmt"

	"github.com/d3sw/floop/types"
	//"github.com/nats-io/go-nats"
	"github.com/nats-io/go-nats-streaming"
)

// NatsStreamdHandler is handler to publish lifecycle events to gnatsd
type NatsStreamdHandler struct {
	conf *types.HandlerConfig
	// Nats connection
	conn stan.Conn
}

// Init initializes the connection to the gnatsd cluster
func (lc *NatsStreamdHandler) Init(conf *types.HandlerConfig) error {
	lc.conf = conf

	clusterID, ok := lc.conf.Options.GetString("cluster_id")
	if !ok || clusterID == "" {
		return errors.New("cluster_id required")
	}
	clientID, ok := lc.conf.Options.GetString("client_id")
	if !ok || clientID == "" {
		return errors.New("client_id required")
	}

	conn, err := stan.Connect(clusterID, clientID)
	if err == nil {
		lc.conn = conn
	}
	return err
}

// Handle publishes to gnatsd.  The config is the normalized
// config built using data from the child process.  This may be different from the one
// used in Init
func (lc *NatsStreamdHandler) Handle(event *types.Event, conf *types.HandlerConfig) (map[string]interface{}, error) {
	// Get topic from config
	topic, ok := conf.Options.GetString("topic")
	if !ok || topic == "" {
		return nil, fmt.Errorf("topic not specified")
	}

	fmt.Printf("[gnatsd] phase=%s topic=%s %+v\n", event.Type, topic, event.Data)

	// Publish the body as bytes
	err := lc.conn.Publish(topic, []byte(conf.Body))

	return nil, err
}

// CloseConnection closes the nats stream connection
func (lc *NatsStreamdHandler) CloseConnection() error {
	lc.conn.Close()
	return nil
}
