package events

import "fmt"

type InputEvent struct {
	Line string
}

func (e InputEvent) Topic(args ...interface{}) string {
	return fmt.Sprintf("session.%s.input", args[0].(string))
}

func (e InputEvent) Payload(args ...interface{}) interface{} {
	if len(args) > 0 {
		return InputEvent{
			Line: args[0].(string),
		}
	}

	return InputEvent{}
}
