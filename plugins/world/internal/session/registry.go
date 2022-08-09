package session

import (
	"fmt"
	"os"
	"strings"

	"github.com/mjolnir-mud/engine/internal/logger"
	"github.com/mjolnir-mud/engine/internal/pubsub"
	"github.com/mjolnir-mud/engine/pkg/events"
	"github.com/rs/zerolog"
)

var sessions map[string]*Session
var log zerolog.Logger
var lineHandlers []func(id string, line string) error
var sessionStartedHandlers []func(id string) error
var sessionStoppedHandlers []func(id string) error
var assertSessionSubscription *pubsub.Subscription
var sessionStoppedSubscription *pubsub.Subscription

func Start() {
	log = logger.Instance.With().Str("service", "sessionRegistry").Logger()
	assertSessionSubscription = pubsub.
		Subscribe(events.AssertSessionTopic, events.AssertSession, handleAssertSession)
	sessionStoppedSubscription = pubsub.
		PSubscribe(events.SessionStoppedTopic("*"), events.SessionStopped, handleSessionStopped)

	err := pubsub.Publish(events.SessionManagerStartedTopic, events.SessionManagerStarted)

	if err != nil {
		log.Error().Err(err).Msg("error session manager started event")
		os.Exit(1)
	}
}

// Stop stops the registry
func Stop() {
	log.Info().Msg("stopping session registry")
	for _, s := range sessions {
		s.Stop()
	}

	sessionStoppedSubscription.Stop()
	assertSessionSubscription.Stop()

	err := pubsub.Publish(events.SessionManagerStoppedTopic, events.SessionManagerStopped)

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

func SendLine(id string, data string) {
	// break up the data by lines
	lines := strings.Split(data, "\n")

	// write each line to the connection
	for _, line := range lines {
		err := pubsub.Publish(events.SendLineTopic(id), events.SendLineEvent{Line: fmt.Sprintf("%s\r\n", line)})

		if err != nil {
			log.Error().Err(err).Msg("error writing to connection")
			break
		}
	}
}

func SendLineF(id string, data string, args ...interface{}) {
	SendLine(id, fmt.Sprintf(data, args...))
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
