package movement

import (
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/character"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/room"
	templates2 "github.com/mjolnir-mud/engine/plugins/world/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

func moveSessionCharacterInDirection(sess session.Session, dir string) error {
	currentLocation := character.GetCurrentCharacterLocationForSession(sess)
	currentCharacter := character.GetCurrentCharacter(sess)
	currentCharacterName := character.GetCharacterName(currentCharacter)

	loc, err := entity_registry.GetStringFromHashComponent(currentLocation, "exits", dir)

	if err != nil {
		switch err.Error() {
		case "redis: nil":
			sess.WriteToConnection(templates.RenderTemplate("walking", &templates2.WalkingContext{
				Direction: dir,
				Focus:     "no-exit",
			}))
			return nil
		default:
			sess.WriteToConnection(err.Error())
			return err
		}
	}

	room.MoveWithMessageForSession(
		sess,
		currentCharacter,
		currentLocation,
		loc,
		templates.RenderTemplate("walking", &templates2.WalkingContext{
			Direction: dir,
			Focus:     "self",
		}),
		templates.RenderTemplate("walking", &templates2.WalkingContext{
			Direction: dir,
			Focus:     "other",
			Name: 	currentCharacterName,
		}),
	)

	return nil
}
