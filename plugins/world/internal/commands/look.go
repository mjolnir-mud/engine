package commands

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/world/internal/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/character"
)

type Look struct{}

func (l *Look) Run(sess session.Session) error {
	loc := character.GetCurrentCharacterLocationForSession(sess)
	n, err := ecs.GetStringComponent(loc, "name")

	if err != nil {
		return err
	}

	return nil
}
