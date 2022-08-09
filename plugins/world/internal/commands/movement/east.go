package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type East struct{}

func (n *East) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "east")
}
