package movement

type Out struct{}

func (n *Out) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "out")
}
