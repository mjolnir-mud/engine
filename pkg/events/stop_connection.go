package events

type StopConnection struct {
	UUID   string
	Reason string
}
