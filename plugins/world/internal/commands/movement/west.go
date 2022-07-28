package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type West struct{}

func (n *West) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "west")
}
