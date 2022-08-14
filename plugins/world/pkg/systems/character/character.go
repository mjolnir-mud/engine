package character

import (
	"github.com/mjolnir-mud/engine/internal/session_registry"

	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
)

type character struct{}

func (e character) Name() string {
	return "expiration"
}

func (e character) Component() string {
	return "type"
}

func (e character) Match(_ string, val interface{}) bool {
	return val.(string) == "character"
}

func (e character) WorldStarted() {}

func (e character) ComponentAdded(entityId string, _ string, _ interface{}) error { return nil }

func (e character) ComponentUpdated(entityId string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e character) ComponentRemoved(entityId string, _ string, _ interface{}) error { return nil }

func (e character) MatchingComponentAdded(entityId string, _ string, _ interface{}) error { return nil }

func (e character) MatchingComponentUpdated(entityId string, _ string, _ interface{}, _ interface{}) error {
	return nil
}

func (e character) MatchingComponentRemoved(entityId string, _ string, _ interface{}) error {
	return nil
}

func GetCurrentCharacter(session reactor.Session) string {
	return session.GetStringFromStore("characterID")
}

func WriteToCharacterConnection(characterId string, message string) {
	sessId, err := entity_registry.GetStringComponent(characterId, "sessionId")

	if err != nil {
		return
	}

	session_registry.Registry.WriteToConnection(sessId, message)
}

func GetCurrentCharacterLocationForSession(session reactor.Session) string {
	characterId := GetCurrentCharacter(session)

	if characterId == "" {
		return ""
	}

	l, err := entity_registry.GetStringComponent(characterId, "location")

	if err != nil {
		return ""
	}

	return l
}

func GetCharacterName(characterId string) string {
	name, err := entity_registry.GetStringComponent(characterId, "name")

	if err != nil {
		return ""
	}

	return name
}

var System = character{}
