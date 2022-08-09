package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type West struct{}

func (n *West) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "west")
}
