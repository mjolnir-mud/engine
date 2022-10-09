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

package engine

import (
	engineEvents "github.com/mjolnir-mud/engine/events"
	"github.com/mjolnir-mud/engine/uid"
	"github.com/rs/zerolog"
)

type sessionHandler struct {
	session  *Session
	logger   zerolog.Logger
	engine   *Engine
	registry *sessionRegistry
}

// Start starts the session handler.
func (h *sessionHandler) Start() {
	h.logger.Info().Msg("starting session handler")
	err := h.engine.Publish(engineEvents.SessionStartedEvent{
		Id: h.session.Id,
	})

	if err != nil {
		h.registry.remove(h.session)
	}
}

// Stop stops the session handler.
func (h *sessionHandler) Stop() {
	h.logger.Info().Msg("stopping session registry")
	h.registry.remove(h.session)
}

type sessionRegistry struct {
	logger                   zerolog.Logger
	sessions                 map[*uid.UID]*sessionHandler
	engine                   *Engine
	sessionStartSubscription *uid.UID
}

func newSessionRegistry(engine *Engine) *sessionRegistry {
	return &sessionRegistry{
		engine:   engine,
		sessions: make(map[*uid.UID]*sessionHandler),
		logger:   engine.Logger.With().Str("component", "session_registry").Logger(),
	}
}

// Start starts the session registry.
func (r *sessionRegistry) Start() {
	r.logger.Info().Msg("starting session registry")

	r.engine.Subscribe(engineEvents.SessionStartEvent{}, r.handleSessionStartEvent)
}

// Stop stops the session registry.
func (r *sessionRegistry) Stop() {
	r.logger.Info().Msg("stopping session registry")

	r.engine.Unsubscribe(r.sessionStartSubscription)
}

func (r *sessionRegistry) add(session *Session) {
	handler := &sessionHandler{
		session: session,
		logger:  r.logger.With().Str("sessionId", session.Id.String()).Logger(),
		engine:  r.engine,
	}

	r.sessions[session.Id] = handler
	handler.Start()
}

func (r *sessionRegistry) remove(session *Session) {
	handler := r.sessions[session.Id]

	if handler == nil {
		return
	}

	handler.Stop()

	delete(r.sessions, session.Id)
}

func (r *sessionRegistry) handleSessionStartEvent(event EventMessage) {
	sessionStartEvent := &engineEvents.SessionStartedEvent{}

	if err := event.Unmarshal(sessionStartEvent); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal session start event")
		return
	}

	r.logger.Debug().Str("sessionId", sessionStartEvent.Id.String()).Msg("session started")
	session := &Session{
		Id: sessionStartEvent.Id,
	}

	err := r.engine.AddEntityWithId(sessionStartEvent.Id, session)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to add session to entity registry")
		return
	}

	r.add(session)
}