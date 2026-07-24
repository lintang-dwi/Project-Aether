package eventbus

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a system event broadcast over the EventBus.
type Event struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// NewEvent creates a new Event with a unique ID and current timestamp.
func NewEvent(topic, source string, payload map[string]interface{}) Event {
	if payload == nil {
		payload = make(map[string]interface{})
	}
	return Event{
		ID:        uuid.New().String(),
		Topic:     topic,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

// EventHandler is a callback function for processing events.
type EventHandler func(event Event)
