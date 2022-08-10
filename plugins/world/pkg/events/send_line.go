package events

import "fmt"

type SendLineEvent struct {
	Line string
}

func SendLineTopic(uuid string) string {
	return fmt.Sprintf("session.%s.send_line", uuid)
}

func SendLine() interface{} {
	return &SendLineEvent{}
}
