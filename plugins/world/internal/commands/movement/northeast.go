package movement

type Northeast struct{}

func (n *Northeast) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "northeast")
}
