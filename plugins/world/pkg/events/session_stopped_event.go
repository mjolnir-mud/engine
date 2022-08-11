package events

import "fmt"

type SessionStoppedEvent struct {
	UUID string
}

func (e SessionStoppedEvent) Topic(args ...interface{}) string {
	return fmt.Sprintf("session.%s.stopped", args[0].(string))
}

func (e SessionStoppedEvent) Payload(args ...interface{}) interface{} {
	if len(args) > 0 {
		return &SessionStoppedEvent{
			UUID: args[0].(string),
		}
	}

	return &SessionStoppedEvent{}
}
