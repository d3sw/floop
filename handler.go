package floop

import "fmt"

// EventType is a type of event representing a stage in the lifecycle
type EventType string

const (
	EventTypeBegin     EventType = "begin"
	EventTypeProgress  EventType = "progress"
	EventTypeCompleted EventType = "completed"
	EventTypeFailed    EventType = "failed"
)

// Event is a single event in a given lifecycle.  Meta is the user passed in metadata.  The type
// of data will be dependent on the event type.
type Event struct {
	Type EventType              `json:"type"`
	Meta map[string]interface{} `json:"meta"`
	Data interface{}            `json:"data"`
}

type Handler interface {
	Handle(event *Event) (map[string]interface{}, error)
}

// EchoHandler implements a Handler that simply echoes back the input
type EchoHandler struct {
}

// Handle echos back input data before process starts
func (lc *EchoHandler) Handle(event *Event) (map[string]interface{}, error) {
	fmt.Printf("[Echo] %+v\n", event)
	return nil, nil
}