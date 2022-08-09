package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Up struct{}

func (n *Up) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "up")
}
