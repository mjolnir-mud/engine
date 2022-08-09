package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Out struct{}

func (n *Out) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "out")
}
