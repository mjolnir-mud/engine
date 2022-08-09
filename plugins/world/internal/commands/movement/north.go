package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type North struct{}

func (n *North) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "north")
}
