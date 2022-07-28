package commands

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/character"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type Look struct{}

func (l *Look) Run(sess session.Session) error {
	loc := character.GetCurrentCharacterLocationForSession(sess)
	n, err := entity_registry.GetStringComponent(loc, "name")

	if err != nil {
		return err
	}

	sess.WriteToConnection(n)

	return nil
}
