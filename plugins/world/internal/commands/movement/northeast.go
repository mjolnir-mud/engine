package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Northeast struct{}

func (n *Northeast) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "northeast")
}
