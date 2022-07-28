package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Northwest struct{}

func (n *Northwest) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "northwest")
}
