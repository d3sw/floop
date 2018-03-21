package types

// EventType is a type of event representing a stage in the lifecycle
type EventType string

const (
	EventTypeBegin     EventType = "begin"
	EventTypeProgress  EventType = "progress"
	EventTypeCompleted EventType = "completed"
	EventTypeFailed    EventType = "failed"
	EventTypeCanceled  EventType = "canceled"
)

// Event is a single event in a given lifecycle.  Meta is the user passed in metadata.  The type
// of data will be dependent on the event type.
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Meta      map[string]interface{} `json:"meta"`
	Data      interface{}            `json:"data"`
}
