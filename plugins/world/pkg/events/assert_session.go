package events

type AssertSessionEvent struct {
	UUID        string
	ConnectedAt int64
	LastInputAt int64
	RemoteAddr  string
}

const AssertSessionTopic = "world.session.assert"

func AssertSession() interface{} {
	return &AssertSessionEvent{}
}
