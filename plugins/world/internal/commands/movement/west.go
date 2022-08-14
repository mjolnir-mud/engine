package movement

type West struct{}

func (n *West) Run(id string) error {
	return moveSessionCharacterInDirection(id, "west")
}
