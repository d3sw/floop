package handlers

import (
	"testing"
	"time"

	"github.com/d3sw/floop/types"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

func Test_NatsStreamHandler(t *testing.T) {

	conf := &types.HandlerConfig{
		URI: nats.DefaultURL,
		Options: types.Options{
			"client_id":  "clientID",
			"cluster_id": "clusterID",
			"topic":      "test",
		},
		Body: "foobar",
	}

	h := &NatsStreamdHandler{}
	err := h.Init(conf)
	if err != nil {
		t.Fatal(err)
	}

	event := &types.Event{
		Type:      types.EventTypeBegin,
		Timestamp: time.Now().UnixNano(),
	}
	_, err = h.Handle(event, conf)
	assert.Nil(t, err)
	if err != nil {
		t.Error(err)
	}
}
