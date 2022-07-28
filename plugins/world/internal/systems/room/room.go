package room

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/character"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/location"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type room struct{}

func (e room) Name() string {
	return "expiration"
}

func (e room) Component() string {
	return "type"
}

func (e room) Match(_ string, value interface{}) bool {
	return value.(string) == "room"
}

func (e room) WorldStarted() {}

func (e room) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e room) ComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error { return nil }

func (e room) ComponentRemoved(_ string, _ string, _ interface{}) error { return nil }

func (e room) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e room) MatchingComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e room) MatchingComponentRemoved(_ string, _ string, _ interface{}) error { return nil }

func Move(entityID string, _ string, toID string) {
	location.Set(entityID, toID)
}

func MoveWithMessageForSession(session session.Session, entityID string, fromID string, toID string, entityMessage string, othersMessage string) {
	characterEntitiesAtLocation := location.AtLocationByType(fromID, "character")
	createCharacter := character.GetCurrentCharacter(session)

	if othersMessage != "" {
		// for each of the characters at the location, send the otherMessage
		for _, characterEntity := range characterEntitiesAtLocation {
			if characterEntity != createCharacter {
				character.WriteToCharacterConnection(characterEntity, othersMessage)
			}
		}
	}

	if entityMessage != "" {
		// write the entityMessage to the character
		session.WriteToConnection(entityMessage)
	}

	// Actually move the entity
	Move(entityID, fromID, toID)
}

var System = room{}
