package movement

type Southeast struct{}

func (n *Southeast) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "southeast")
}
