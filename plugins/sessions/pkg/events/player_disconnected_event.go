package events

import "fmt"

type PlayerDisconnectedEvent struct {
	Id string
}

func (e PlayerDisconnectedEvent) Topic(args ...interface{}) string {
	return fmt.Sprintf("session.%s.stopped", args[0].(string))
}

func (e PlayerDisconnectedEvent) Payload(args ...interface{}) interface{} {
	if len(args) > 0 {
		return &PlayerDisconnectedEvent{
			Id: args[0].(string),
		}
	}

	return &PlayerDisconnectedEvent{}
}
