package session

import (
	"github.com/mjolnir-mud/engine/internal/pubsub"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/rs/zerolog"
)

// Session represents the link between the world and portal.
type Session struct {
	id                      string
	logger                  zerolog.Logger
	receiveLineSubscription pubsub.Subscription
}

func New(e *events.AssertSessionEvent) *Session {
	return &Session{
		id:     e.UUID,
		logger: logger.Instance.With().Str("session", e.UUID).Logger(),
	}
}

func (s Session) Start() {
	s.logger.Info().Msg("starting")
	for _, h := range sessionStartedHandlers {
		err := h(s.id)

		if err != nil {
			s.logger.Error().Err(err).Msg("error starting session")
			s.Stop()
		}
	}

	pubsub.Subscribe(events.InputEventTopic(s.id), events.Input, func(payload interface{}) {
		event, ok := payload.(*events.InputEvent)

		if !ok {
			s.logger.Error().Msg("error casting event to SendToWorldEvent")
			return
		}

		for _, handler := range lineHandlers {
			err := handler(s.id, event.Line)

			if err != nil {
				s.logger.Error().Err(err).Msg("error handling line")
				s.Stop()
				return
			}
		}

	})
}

func (s Session) Stop() {
	s.logger.Info().Msg("stopping")
	for _, h := range sessionStoppedHandlers {
		err := h(s.id)

		if err != nil {
			s.logger.Error().Err(err).Msg("error stopping session")
		}
	}

	err := pubsub.Publish(events.SessionStoppedTopic(s.id), events.SessionStoppedEvent{
		UUID: s.id,
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("error publishing session stopped event")
	}
}
