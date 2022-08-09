package events

type SessionManagerStartedEvent struct {
	Time int64
}

const SessionManagerStartedTopic = "session_manager.started"

func SessionManagerStarted() interface{} {
	return &SessionManagerStartedEvent{}
}
