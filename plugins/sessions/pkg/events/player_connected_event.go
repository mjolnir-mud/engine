package events

type PlayerConnectedEvent struct {
	Id string
}

func (e PlayerConnectedEvent) Topic(_ ...interface{}) string {
	return "player_connected"
}

func (e PlayerConnectedEvent) Payload(args ...interface{}) interface{} {
	if len(args) == 0 {
		return &PlayerConnectedEvent{}
	}

	return &PlayerConnectedEvent{
		Id: args[0].(string),
	}
}
