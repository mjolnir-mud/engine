package movement

import (
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
)

type South struct{}

func (n *South) Run(sess session.Session) error {
	return moveSessionCharacterInDirection(sess, "south")
}
