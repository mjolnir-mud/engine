package movement

import "github.com/mjolnir-mud/engine/plugins/world/pkg/session"

type Down struct{}

func (n *Down) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "down")
}
