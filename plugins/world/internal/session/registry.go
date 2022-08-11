package session

import (
	"fmt"
	"strings"

	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/rs/zerolog"
)

var sessions map[string]*Session
var log zerolog.Logger
var lineHandlers []func(id string, line string) error
var sessionStartedHandlers []func(id string) error
var sessionStoppedHandlers []func(id string) error
var assertSessionSubscription *redis.Subscription
var sessionStoppedSubscription *redis.Subscription

func RegisterLineHandler(f func(id string, line string) error) {
	lineHandlers = append(lineHandlers, f)
}

func RegisterSessionStartedHandler(f func(id string) error) {
	sessionStartedHandlers = append(sessionStartedHandlers, f)
}

func RegisterSessionStoppedHandler(f func(id string) error) {
	sessionStoppedHandlers = append(sessionStoppedHandlers, f)
}

func SendLine(id string, data string) {
	// break up the data by lines
	lines := strings.Split(data, "\n")

	// write each line to the connection
	for _, line := range lines {
		err := redis.Publish(events.SendLineEvent{}, id, line)

		if err != nil {
			log.Error().Err(err).Msg("error writing to connection")
			break
		}
	}
}

func SendLineF(id string, data string, args ...interface{}) {
	SendLine(id, fmt.Sprintf(data, args...))
}

func StartRegistry() {
	log = logger.Instance.With().Str("service", "sessionRegistry").Logger()
	assertSessionSubscription = redis.Subscribe(events.AssertSessionEvent{}, handleAssertSession)
	sessionStoppedSubscription = redis.PSubscribe(events.SessionStoppedEvent{}, "*", handleSessionStopped)

	err := redis.Publish(events.SessionRegistryStartedEvent{})

	if err != nil {
		log.Error().Err(err).Msg("error session manager started event")
		panic(err)
	}
}

// StopRegistry stops the registry
func StopRegistry() {
	log.Info().Msg("stopping session registry")
	for _, s := range sessions {
		s.Stop()
	}

	sessionStoppedSubscription.Stop()
	assertSessionSubscription.Stop()

	err := redis.Publish(events.SessionRegistryStoppedEvent{})

	if err != nil {
		log.Error().Err(err).Msg("error session manager stopped event")
	}
}

func StopSession(uuid string) {
	sess, ok := sessions[uuid]

	if !ok {
		return
	}

	sess.Stop()
	remove(sess)
}

func add(s *Session) {
	sessions[s.id] = s
	s.Start()
}

func remove(s *Session) {
	delete(sessions, s.id)
}

func handleAssertSession(payload interface{}) {
	log.Debug().Msg("handling assert session event")
	assertSession := payload.(events.AssertSessionEvent)

	sess := New(&assertSession)
	add(sess)
}

func handleSessionStopped(payload interface{}) {
	log.Debug().Msg("handling session stopped event")
	sessionStopped := payload.(events.SessionStoppedEvent)

	remove(sessions[sessionStopped.UUID])
}
