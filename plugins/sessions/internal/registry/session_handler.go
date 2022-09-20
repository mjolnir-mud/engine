/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/plugins/ecs"
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
	return engine.Publish(events.PlayerOutputEvent{Id: h.Id, Line: line + "\r\n"})
}

func (h *sessionHandler) Start() {
	h.logger.Debug().Msg("starting")
	err := ecs.AddEntityWithID("session", h.Id, map[string]interface{}{})

	if err != nil {
		h.logger.Error().Err(err).Msg("error creating session entity")
		return
	}

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

func (h *sessionHandler) receiveLine(line string) {
	for _, handler := range receiveLineHandlers {
		err := handler(h.Id, line)

		if err != nil {
			h.logger.Error().Err(err).Msg("error handling line")
			h.Stop()
			return
		}
	}
}

func (h *sessionHandler) sendLine(line string) {
	for _, handler := range sendLineHandlers {
		err := handler(h.Id, line)

		if err != nil {
			h.logger.Error().Err(err).Msg("error handling line")
			h.Stop()

			return
		}
	}
}
