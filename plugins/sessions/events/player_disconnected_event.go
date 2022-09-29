package events

import "fmt"

type PlayerDisconnectedEvent struct {
	Id string
}

func (e PlayerDisconnectedEvent) Topic() string {
	return fmt.Sprintf("%s.disconnected", e.Id)
}
