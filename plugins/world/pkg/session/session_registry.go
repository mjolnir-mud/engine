package session

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/rs/zerolog"
)

type registry struct {
	sessions map[string]*Session
	stop     chan bool
	logger   zerolog.Logger
}

// Start starts the registry
func (r *registry) Start() {
	r.logger.Info().Msg("starting session registry")

	// set up assert session subscriber
	_, err := engine.SubscribeToEvent(events.AssertSessionTopic(), r.handleAssertSession)

	// error out and exit if we can't rake the subscription
	if err != nil {
		r.logger.Error().Err(err).Msg("error subscribing to event")
		os.Exit(1)
	}

	err = engine.PublishEvent(events.SessionManagerStartedTopic(), events.SessionManagerStarted{
		Time: time.Now().UTC().Unix(),
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("error publishing event")
		os.Exit(1)
	}

	// wait for the stop signal
	<-r.stop
}

func (r *registry) WriteToConnection(id string, data string) {
	// break up the data by lines
	lines := strings.Split(data, "\n")

	// write each line to the connection
	for _, line := range lines {
		err := engine.PublishEvent(
			events.WriteToConnectionTopic(id), events.WriteToConnection{Line: fmt.Sprintf("%s\r\n", line)},
		)

		if err != nil {
			r.logger.Error().Err(err).Msg("error writing to connection")
			break
		}
	}
}

func (r *registry) WriteToConnectionF(id string, data string, args ...interface{}) {
	r.WriteToConnection(id, fmt.Sprintf(data, args...))
}

// Stop stops the registry
func (r *registry) Stop() {
	r.logger.Info().Msg("stopping session registry")

	// signal the stop channel
	r.stop <- true
}

func (r *registry) add(s *Session) {
	r.sessions[s.id] = s
	s.Start()
}

func (r *registry) remove(s *Session) {
	delete(r.sessions, s.id)
}

func (r *registry) handleAssertSession(as *events.AssertSession) {
	r.logger.Debug().Msg("handling assert session event")

	sess := New(as)
	r.add(sess)
}

var Registry = &registry{
	sessions: make(map[string]*Session),
	stop:     make(chan bool),
	logger:   logger.Logger.With().Str("service", "sessionManager").Logger(),
}
