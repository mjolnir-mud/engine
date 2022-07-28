package movement

import (
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type North struct{}

func (n *North) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "north")
}
