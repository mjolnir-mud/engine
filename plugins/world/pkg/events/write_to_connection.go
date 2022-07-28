package events

import "fmt"

type WriteToConnection struct {
	Line string `json:"line"`
}

func WriteToConnectionTopic(uuid string) string {
	return fmt.Sprintf("world.session.%s.write_to_connection", uuid)
}
