package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Out struct{}

func (n *Out) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "out")
}
