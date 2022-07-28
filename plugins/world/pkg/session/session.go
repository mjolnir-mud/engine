package session

import (
	"context"
	"fmt"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// Session represents the link between the world and portal.
type Session struct {
	id             string
	entity         string
	logger         zerolog.Logger
	sendToWorldSub *nats.Subscription
}

func New(e *events.AssertSession) *Session {
	nsa := map[string]interface{}{
		"id":          e.UUID,
		"remoteAddr":  e.RemoteAddr,
		"connectedAt": e.ConnectedAt,
		"lastInputAt": e.LastInputAt,
	}

	err := entity_registry.Add("session", e.UUID, nsa)

	if err != nil {
		return nil
	}

	return &Session{
		id:     e.UUID,
		logger: logger.Logger.With().Str("session_id", e.UUID).Logger(),
		entity: e.UUID,
	}
}

func (s Session) Start() {
	s.logger.Info().Msg("starting")

	sub, err := engine.SubscribeToEvent(events.SendToWorldTopic(s.id), func(e *events.SendToWorld) {
		c, err := s.GetController()

		if err != nil {
			s.logger.Error().Err(err).Msg("error getting controller")
			s.Stop()
			return
		}

		err = c.HandleInput(s, e.Line)

		if err != nil {
			s.logger.Error().Err(err).Msg("error handling input")
			s.Stop()
			return
		}
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("error subscribing to event")
		s.Stop()
		return
	}

	s.sendToWorldSub = sub

	c, err := s.GetController()

	if err != nil {
		s.logger.Error().Err(err).Msg("error getting controller")
		s.Stop()
		return
	}

	err = c.Start(s)

	if err != nil {
		s.logger.Error().Err(err).Msg("error starting controller")
		s.Stop()
		return
	}

}

func (s Session) Stop() {
	s.logger.Info().Msg("stopping")
}

func (s Session) ID() string {
	return s.id
}

func (s Session) Resume() {
}

func (s Session) SetController(name string) {
	s.logger.Debug().Msgf("setting controller to %s", name)
	entity_registry.RemoveComponent(s.entity, "flash")
	s.SetInStore("controller", name)

	c, err := s.GetController()

	if err != nil {
		s.logger.Error().Err(err).Msg("error starting controller")
		s.Stop()
		return
	}

	err = c.Start(s)

	if err != nil {
		s.logger.Error().Err(err).Msg("error starting controller")
		s.Stop()
	}

	s.logger.Debug().Msgf("controller %s started", name)
}

// GetController returns the current controller or an error if a controller of that controller does not exist
func (s Session) GetController() (Controller, error) {
	name := s.GetStringFromStore("controller")

	c := ControllerRegistry.Get(name)

	if c == nil {
		return nil, fmt.Errorf("controller %s does not exist", name)
	}

	return c, nil
}

func (s Session) SetInStore(key string, value interface{}) {
	entity_registry.SetInHashComponent(s.entity, "store", key, value)

}

func (s Session) SetStringInStore(key string, value string) {
	s.SetInStore(key, value)
}

func (s Session) SetInFlash(key string, value interface{}) {
	entity_registry.SetInHashComponent(s.entity, "flash", key, value)
}

// WriteToConnection writes the given string to the connection. It breaks it up by lines and sends each line to the
// connection. If an error occurs, it will stop the session.
func (s Session) WriteToConnection(data string) {
	Registry.WriteToConnection(s.id, data)
}

func (s Session) WriteToConnectionF(format string, args ...interface{}) {
	Registry.WriteToConnectionF(s.id, format, args...)
}

func (s Session) GetStringFromStore(key string) string {
	str, err := entity_registry.GetStringFromHashComponent(s.entity, "store", key)

	if err != nil {
		s.logger.Error().Err(err).Msgf("error getting string from store %s", key)
		s.Stop()
	}

	return str
}

func (s Session) GetIntFromStore(key string) int {
	i, err := entity_registry.GetIntFromHashComponent(s.entity, "store", key)

	if err != nil {
		s.logger.Error().Err(err).Msgf("error getting int from store %s", key)
		s.Stop()
	}

	return i
}

func (s Session) GetStringFromFlash(key string) string {
	str, err := entity_registry.GetStringFromHashComponent(s.entity, "flash", key)

	if err != nil {
		s.logger.Error().Err(err).Msgf("error getting string from flash %s", key)
		s.Stop()
	}

	return str
}

func (s Session) GetIntFromFlash(key string) int {
	i, err := entity_registry.GetIntFromHashComponent(s.entity, "flash", key)

	if err != nil {
		s.logger.Error().Err(err).Msgf("error getting int from flash %s", key)
		s.Stop()
	}

	return i
}

func (s Session) Logger() *zerolog.Logger {
	return &s.logger
}

func (s Session) Context() context.Context {
	return context.WithValue(context.Background(), "session", s)
}
