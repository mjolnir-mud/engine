package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type East struct{}

func (n *East) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "east")
}
