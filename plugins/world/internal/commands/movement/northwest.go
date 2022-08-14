package movement

type Northwest struct{}

func (n *Northwest) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "northwest")
}
