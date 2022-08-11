package event

type Event interface {
	// Topic returns the topic the event should be published to. It accepts the same arguments as the Payload method,
	// and should return a string representation of the topic.
	Topic(args ...interface{}) string

	// Payload returns the payload of the event. It accepts a set of arguments that can be used to construct the payload.
	Payload(args ...interface{}) interface{}
}
