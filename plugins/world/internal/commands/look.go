package commands

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/character"
)

type Look struct{}

func (l *Look) Run(sess reactor.Session) error {
	loc := character.GetCurrentCharacterLocationForSession(sess)
	n, err := entity_registry.GetStringComponent(loc, "name")

	if err != nil {
		return err
	}

	sess.WriteToConnection(n)

	return nil
}
