package events

type SessionManagerStoppedEvent struct{}

const SessionManagerStoppedTopic = "session_manager.stopped"

func SessionManagerStopped() interface{} {
	return &SessionStoppedEvent{}
}
