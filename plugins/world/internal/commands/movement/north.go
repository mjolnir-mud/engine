package movement

type North struct{}

func (n *North) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "north")
}
