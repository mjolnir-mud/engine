package events

type SessionRegistryStartedEvent struct{}

func (e SessionRegistryStartedEvent) Topic() string {
	return "session_manager.started"
}
