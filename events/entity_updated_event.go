package events

import (
	"fmt"

	"github.com/mjolnir-engine/engine/uid"
)

type EntityUpdatedEvent struct {
	Id uid.UID
}

func (e EntityUpdatedEvent) Topic() string {
	return fmt.Sprintf("entity.%s.updated", e.Id)
}

func (e EntityUpdatedEvent) AllTopics() string {
	return "entity.*.updated"
}
