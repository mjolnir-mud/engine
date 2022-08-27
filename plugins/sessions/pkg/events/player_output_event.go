package events

import "fmt"

type PlayerOutputEvent struct {
	Id   string
	Line string
}

func (e PlayerOutputEvent) Topic() string {
	return fmt.Sprintf("%s.output", e.Id)
}
