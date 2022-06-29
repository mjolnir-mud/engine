package events

type AssertConnection struct {
	UUID        string
	ConnectedAt int64
	LastInputAt int64
	RemoteAddr  string
}

func AssertConnectionTopic() string {
	return "engine.assert_connection"
}
