package events

type SessionRegistryStoppedEvent struct{}

func (e SessionRegistryStoppedEvent) Topic(_ ...interface{}) string {
	return "session_manager.stopped"
}

func (e SessionRegistryStoppedEvent) Payload(_ ...interface{}) interface{} {
	return SessionRegistryStoppedEvent{}
}
