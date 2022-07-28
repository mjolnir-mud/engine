package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Up struct{}

func (n *Up) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "up")
}
