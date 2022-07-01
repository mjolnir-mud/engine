package engine

import "encoding/json"

type EventMessage struct {
	Topic   string
	payload string
}

// NewEventMessage creates a new EventMessage.
func NewEventMessage(topic string, payload string) *EventMessage {
	return &EventMessage{
		Topic:   topic,
		payload: payload,
	}
}

// Unmarshal parses the payload of the event message from JSON,
func (e *EventMessage) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(e.payload), v)
}
