package commands

import (
	"strings"

	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/location"
	templates2 "github.com/mjolnir-mud/engine/plugins/world/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type Say struct {
	Text []string `arg:"required"`
}

func (c *Say) Run(sess session.Session) error {
	characterId := sess.GetStringFromStore("characterID")
	roomId, err := entity_registry.GetStringComponent(characterId, "location")

	if err != nil {
		return err
	}

	name, err := entity_registry.GetStringComponent(characterId, "name")

	if err != nil {
		return err
	}

	entityIds := location.AtLocation(roomId)

	sess.WriteToConnection(templates.RenderTemplate("say", &templates2.SayContext{
		Focus: "self",
		Message: strings.Join(c.Text, " "),
	}))

	for _, entityId := range entityIds {
		sessId, err := entity_registry.GetStringComponent(entityId, "sessionID")

		if err != nil {
			continue
		}

		if sessId == sess.ID() {
			continue
		}

		session.Registry.WriteToConnection(sessId, templates.RenderTemplate("say", &templates2.SayContext{
			Focus: "other",
			Message: strings.Join(c.Text, " "),
			Name: name,
		}))
	}

	return nil
}

