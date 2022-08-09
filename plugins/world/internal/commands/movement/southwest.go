package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Southwest struct{}

func (n *Southwest) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "southwest")
}
