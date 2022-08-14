package movement

type Up struct{}

func (n *Up) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "up")
}
