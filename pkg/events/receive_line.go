package events

import "fmt"

// ReceiveLine is an event to send a line to the world.
type InputEvent struct {
	Line string
}

func InputEventTopic(uuid string) string {
	return fmt.Sprintf("session.%s.input", uuid)
}

func Input() interface{} {
	return &InputEvent{}
}
