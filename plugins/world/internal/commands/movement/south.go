package movement

type South struct{}

func (n *South) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "south")
}
