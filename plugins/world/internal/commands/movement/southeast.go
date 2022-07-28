package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Southeast struct{}

func (n *Southeast) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "southeast")
}
