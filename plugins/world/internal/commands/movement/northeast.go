package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Northeast struct{}

func (n *Northeast) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "northeast")
}
