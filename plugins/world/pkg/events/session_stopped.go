package events

import "fmt"

type SessionStoppedEvent struct {
	UUID string
}

func SessionStoppedTopic(uuid string) string {
	return fmt.Sprintf("session.%s.stopped", uuid)
}

func SessionStopped() interface{} {
	return &SessionStoppedEvent{}
}
