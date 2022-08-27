package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/plugins/sessions/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"

	"github.com/rs/zerolog"
)

var sessionHandlers map[string]*sessionHandler
var playerConnectedSubscription engine.Subscription
var lineHandlers []func(id string, line string) error
var sessionStartedHandlers []func(id string) error
var sessionStoppedHandlers []func(id string) error
var log zerolog.Logger

func RegisterLineHandler(f func(id string, line string) error) {
	lineHandlers = append(lineHandlers, f)
}

func RegisterSessionStartedHandler(f func(id string) error) {
	sessionStartedHandlers = append(sessionStartedHandlers, f)
}

func RegisterSessionStoppedHandler(f func(id string) error) {
	sessionStoppedHandlers = append(sessionStoppedHandlers, f)
}

func Start() {
	log = logger.Instance.With().Str("service", "sessionRegistry").Logger()
	sessionHandlers = make(map[string]*sessionHandler, 0)
	playerConnectedSubscription = engine.Subscribe(events.PlayerConnectedEvent{}, handlePlayerConnected)
	lineHandlers = make([]func(id string, line string) error, 0)
	sessionStartedHandlers = make([]func(id string) error, 0)
	sessionStoppedHandlers = make([]func(id string) error, 0)

	err := engine.Publish(events.SessionRegistryStartedEvent{})

	if err != nil {
		log.Error().Err(err).Msg("error session manager started event")
		panic(err)
	}
}

func Stop() {
	log.Info().Msg("stopping")
	playerConnectedSubscription.Stop()

	for _, s := range sessionHandlers {
		s.Stop()
	}
}

func StopSession(id string) {
	log.Debug().Str("id", id).Msg("stopping session")
	sess := sessionHandlers[id]
	sess.Stop()
}

func add(s *sessionHandler) {
	sessionHandlers[s.Id] = s
	s.Start()
}

func remove(id string) {
	delete(sessionHandlers, id)
}

func handlePlayerConnected(payload event.EventPayload) {
	log.Debug().Msg("handling player connected event")
	newConnection := &events.PlayerConnectedEvent{}
	err := payload.Unmarshal(newConnection)

	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling player connected event")
		return
	}

	sess := NewSessionHandler(newConnection.Id)
	add(sess)

	sess.Start()
}
