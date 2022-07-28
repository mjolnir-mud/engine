package events

import "fmt"

// SendToWorld is an event to send a line to the world.
type SendToWorld struct {
	Line string `json:"line"`
}

func SendToWorldTopic(uuid string) string {
	return fmt.Sprintf("world.session.%s.send_to_world", uuid)
}
