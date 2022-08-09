package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type South struct{}

func (n *South) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "south")
}
