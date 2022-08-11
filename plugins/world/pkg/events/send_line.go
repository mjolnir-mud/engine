package events

import "fmt"

type SendLineEvent struct {
	Line string
}

func (e SendLineEvent) Topic(args ...interface{}) string {
	return fmt.Sprintf("session.%s.output", args[0].(string))
}

func (e SendLineEvent) Payload(args ...interface{}) interface{} {
	if len(args) > 0 {
		return SendLineEvent{
			Line: args[1].(string),
		}
	}

	return SendLineEvent{}
}
