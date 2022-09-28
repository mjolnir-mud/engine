package event

import (
	"encoding/json"
)

type Event interface {
	// Topic returns the topic the event should be published to. It accepts the same arguments as the Payload method,
	// and should return a string representation of the topic.
	Topic() string
}

// EventPayload is the payload of an event. Call `Unmarshal` on it to unmarshal the payload into a struct.
type EventPayload struct {
	Payload []byte
}

// Unmarshal unmarshals the payload into the given data_source.
func (ep *EventPayload) Unmarshal(i interface{}) error {
	return json.Unmarshal(ep.Payload, i)
}
