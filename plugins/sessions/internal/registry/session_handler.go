package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/plugins/sessions/pkg/events"
	"github.com/rs/zerolog"
)

type sessionHandler struct {
	Id                       string
	lineSubscription         engine.Subscription
	disconnectedSubscription engine.Subscription
	logger                   zerolog.Logger
}

func NewSessionHandler(id string) *sessionHandler {
	ds := engine.Subscribe(events.PlayerDisconnectedEvent{Id: id}, func(e event.EventPayload) {
		StopSession(id)
	})

	s := &sessionHandler{
		Id: id,
		logger: log.
			With().
			Str("service", "sessionHandler").
			Str("id", id).
			Logger(),
		disconnectedSubscription: ds,
	}

	ls := engine.Subscribe(events.PlayerInputEvent{Id: id}, func(payload event.EventPayload) {
		e := &events.PlayerInputEvent{}

		err := payload.Unmarshal(e)

		if err != nil {
			s.logger.Error().Err(err).Msg("error unmarshalling event")
			s.Stop()
			return
		}

		s.receiveLine(e.Line)
	})

	s.lineSubscription = ls

	return s
}

func (h *sessionHandler) SendLine(line string) error {
	return h.sendLine(line)
}

func (h *sessionHandler) Start() {
	h.logger.Debug().Msg("starting")
	for _, handler := range sessionStartedHandlers {
		err := handler(h.Id)

		if err != nil {
			h.logger.Error().Err(err).Msg("error starting session")
			h.Stop()
			return
		}
	}
}

func (h *sessionHandler) Stop() {
	for _, handler := range sessionStoppedHandlers {
		err := handler(h.Id)

		if err != nil {
			h.logger.Error().Err(err).Msg("error stopping session")
			return
		}
	}

	h.disconnectedSubscription.Stop()
	h.lineSubscription.Stop()
	remove(h.Id)
}

func (h *sessionHandler) sendLine(line string) error {
	return engine.Publish(events.PlayerOutputEvent{Id: h.Id, Line: line})
}

func (h *sessionHandler) receiveLine(line string) {
	for _, handler := range lineHandlers {
		err := handler(h.Id, line)

		if err != nil {
			h.logger.Error().Err(err).Msg("error handling line")
			h.Stop()
			return
		}
	}
}
