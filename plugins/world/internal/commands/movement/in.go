package movement

type In struct{}

func (n *In) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "in")
}
