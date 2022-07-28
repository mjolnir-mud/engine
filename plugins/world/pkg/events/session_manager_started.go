package events

type SessionManagerStarted struct {
	Time int64
}

func SessionManagerStartedTopic() string {
	return "world.session_manager.started"
}
