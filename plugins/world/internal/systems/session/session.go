package session

import (
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
)

type session struct{}

func (e session) Name() string {
	return "expiration"
}

func (e session) Match(key string, value interface{}) bool {
	// return false if value is not a string
	if _, ok := value.(string); !ok {
		return false
	}

	// return true if the value is session and the key is _type
	return key == "_type" && value.(string) == "session"
}

func (e session) WorldStarted() {}

func (e session) ComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e session) ComponentUpdated(_ string, _ string, _ interface{}) error { return nil }

func (e session) ComponentRemoved(_ string, _ string, _ interface{}) error { return nil }

func (e session) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }

func (e session) MatchingComponentUpdated(_ string, _ string, _ interface{}) error { return nil }

func (e session) MatchingComponentRemoved(_ string, _ string, _ interface{}) error { return nil }

func GetCommandSets(sess reactor.Session) []string {
	set, err := entity_registry.GetStringsFromSetComponent(sess.ID(), "commandSets")

	if err != nil {
		return []string{}
	}

	return set
}

func AddCommandSet(sess reactor.Session, set string) {
}

func RemoveCommandSet(sess reactor.Session, set string) {
}

var System = session{}
