package movement

type Down struct{}

func (n *Down) Run(sess string) error {
	return moveSessionCharacterInDirection(sess, "down")
}
