package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Southeast struct{}

func (n *Southeast) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "southeast")
}
