package movement

type East struct{}

func (n *East) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "east")
}
