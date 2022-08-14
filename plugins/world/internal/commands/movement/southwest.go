package movement

type Southwest struct{}

func (n *Southwest) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "southwest")
}
