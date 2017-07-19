package handlers

import (
	"testing"
	"time"

	"github.com/d3sw/floop/types"
	"github.com/nats-io/go-nats"
)

func TestGnatsdHandler(t *testing.T) {

	conf := &types.HandlerConfig{
		URI:     nats.DefaultURL,
		Options: types.Options{"topic": "test"},
		Body:    "foobar",
	}

	h := &GnatsdHandler{}
	h.Init(conf)
	if _, err := h.Handle(&types.Event{Type: types.EventTypeBegin, Timestamp: time.Now().UnixNano()}, conf); err != nil {
		t.Fatal(err)
	}

	// sb, err := h.conn.SubscribeSync("test")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// msg, err := sb.NextMsg(1 * time.Second)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// t.Logf("%+v", msg)

}
