package movement

import (
	"github.com/mjolnir-mud/engine/pkg/reactor"
)

type Down struct{}

func (n *Down) Run(sess reactor.Session) error {
	return moveSessionCharacterInDirection(sess, "down")
}
