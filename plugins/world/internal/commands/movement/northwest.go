package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Northwest struct{}

func (n *Northwest) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "northwest")
}
