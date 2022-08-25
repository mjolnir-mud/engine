package events

type SessionRegistryStartedEvent struct{}

func (e SessionRegistryStartedEvent) Topic(_ ...interface{}) string {
	return "session_manager.started"
}

func (e SessionRegistryStartedEvent) Payload(_ ...interface{}) interface{} {
	return &SessionRegistryStartedEvent{}

}
