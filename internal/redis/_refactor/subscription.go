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

package _refactor

import (
	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/logger"
	"github.com/rs/zerolog"
)

type Subscription struct {
	pubsub   *redis.PubSub
	stop     chan bool
	callback func(payload engine.EventPayload)
	logger   zerolog.Logger
	event    engine.Event
}

func NewSubscription(e engine.Event, callback func(payload engine.EventPayload)) *Subscription {
	return createSubscription(Subscribe(e.Topic()), e, callback)

}

func NewPatternSubscription(e engine.Event, callback func(payload engine.EventPayload)) *Subscription {
	return createSubscription(PSubscribe(e.Topic()), e, callback)
}

func createSubscription(pubsub *redis.PubSub, e engine.Event, callback func(payload engine.EventPayload)) *Subscription {
	s := &Subscription{
		pubsub:   pubsub,
		stop:     make(chan bool),
		callback: callback,
		event:    e,
		logger: logger.Instance.
			With().
			Str("service", "pubsub").
			Str("topic", pubsub.String()).
			Logger(),
	}

	go func() {
		s.logger.Debug().Msg("starting subscription")
		for {
			select {
			case <-s.stop:
				_ = s.pubsub.Close()
				return
			case msg := <-s.pubsub.Channel():
				if msg == nil {
					s.logger.Debug().Msg("channel closed")
					return
				}

				payloadBytes := []byte(msg.Payload)
				length := len(payloadBytes)

				s.logger.Debug().Msgf("received message: %d", length)

				newEvent := engine.EventPayload{
					Payload: payloadBytes,
				}

				s.callback(newEvent)
			}
		}
	}()

	return s
}

func (s *Subscription) Stop() {
	s.logger.Info().Msg("stopping pubsub")
	s.stop <- true
}
