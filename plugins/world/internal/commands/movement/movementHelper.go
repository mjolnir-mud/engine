package movement

import (
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/templates"
	templates2 "github.com/mjolnir-mud/engine/plugins/world/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/character"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/systems/room"
)

func moveSessionCharacterInDirection(sess string, dir string) error {
	currentLocation := character.GetCurrentCharacterLocationForSession(sess)
	currentCharacter := character.GetCurrentCharacter(sess)
	currentCharacterName := character.GetCharacterName(currentCharacter)

	loc, err := ecs.GetStringFromMapComponent(currentLocation, "exits", dir)

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
			Name:      currentCharacterName,
		}),
	)

	return nil
}
