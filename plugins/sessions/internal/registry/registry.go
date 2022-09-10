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
	"github.com/mjolnir-mud/engine/plugins/controllers/pkg/errors"
	"github.com/mjolnir-mud/engine/plugins/ecs"
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
	log = logger.Instance.With().Str("component", "registry").Logger()
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

func SendLine(id string, line string) error {
	exists, err := ecs.EntityExists(id)

	if err != nil {
		return err
	}

	sess, ok := sessionHandlers[id]

	if !exists && ok {
		StopSession(id)
	}

	if !exists {
		return errors.SessionNotFoundError{SessionId: id}
	}

	return sess.SendLine(line)
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
