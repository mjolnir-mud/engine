package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Southwest struct{}

func (n *Southwest) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "southwest")
}
