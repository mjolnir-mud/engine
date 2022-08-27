package events

import "fmt"

type PlayerInputEvent struct {
	Id   string
	Line string
}

func (e PlayerInputEvent) Topic() string {
	return fmt.Sprintf("%s.input", e.Id)
}
