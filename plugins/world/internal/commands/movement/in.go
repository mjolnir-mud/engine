package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type In struct{}

func (n *In) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "in")
}
