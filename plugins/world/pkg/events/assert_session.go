package events

type AssertSessionEvent struct {
	UUID        string
	ConnectedAt int64
	LastInputAt int64
	RemoteAddr  string
}

func (e AssertSessionEvent) Topic(_ ...interface{}) string {
	return "session_registry.assert_session"
}

func (e AssertSessionEvent) Payload(_ ...interface{}) interface{} {
	return AssertSessionEvent{}
}
