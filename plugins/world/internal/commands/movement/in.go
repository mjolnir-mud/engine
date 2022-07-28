package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type In struct{}

func (n *In) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "in")
}
