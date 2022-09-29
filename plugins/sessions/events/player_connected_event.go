package events

type PlayerConnectedEvent struct {
	Id string
}

func (e PlayerConnectedEvent) Topic() string {
	return "player_connected"
}
