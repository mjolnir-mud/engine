package events

type AssertSession struct {
	UUID        string
	ConnectedAt int64
	LastInputAt int64
	RemoteAddr  string
}

func AssertSessionTopic() string {
	return "world.assert_session"
}
