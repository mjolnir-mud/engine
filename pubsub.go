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
	"context"
	"encoding/json"
	"fmt"
	engineErrors "github.com/mjolnir-mud/engine/errors"
	"github.com/rueian/rueidis"
)

// EventMessage is a message received from a subscription. It can be used to unmarshall the event.
type EventMessage struct {
	message rueidis.PubSubMessage
}

// Unmarshal unmarshalls the event message from JSON into the provided event.
func (e EventMessage) Unmarshal(event interface{}) error {
	err := json.Unmarshal([]byte(e.message.Message), event)

	if err != nil {
		return err
	}

	return nil
}

// Publish publishes events to the event bus. It can publish multiple events at once. It expects a slice of items that
// implement the Event interface. It will marshal the event into JSON and publish it to the topic returned by the
// Topic method. If there is an error publishing the event. This will return a `PublishErrors` error, which contains
// a slice of errors for each event that failed to publish.
func (e *Engine) Publish(events ...Event) error {
	logger := e.Logger.With().Str("component", "publisher").Logger()

	logger.Debug().Int("events", len(events)).Msg("publishing events")
	commands := e.GetPublishCommandsForEvents(events...)

	results := e.redis.DoMulti(
		context.Background(),
		commands...,
	)

	publishErrors := engineErrors.PublishErrors{}

	for _, result := range results {
		if result.Error() != nil {
			publishErrors.Add(result.Error())
		}
	}

	if publishErrors.HasErrors() {
		return &publishErrors
	}

	return nil
}

// GetPublishCommandsForEvents returns a slice of commands that can be used to publish the provided events. This can be
// used with `DoMulti` to publish multiple events along with other Redis commands in a single transaction.
func (e *Engine) GetPublishCommandsForEvents(events ...Event) rueidis.Commands {
	commands := make(rueidis.Commands, len(events))

	for i, event := range events {
		commands[i] = e.redis.B().Publish().Channel(e.topicWithPrefix(event)).Message(rueidis.JSON(event)).Build()
	}

	return commands
}

// Subscribe subscribes an event. The event must implement the `Event` interface. A callback function is to be provided
// which will be called when the event is published, the callback will be passed an `EventMessage` which can be used to
// unmarshall the event.
func (e *Engine) Subscribe(event Event, callback func(event EventMessage)) {
	e.Logger.Debug().Str("topic", e.topicWithPrefix(event)).Msg("subscribing to topic")

	client, cancel := e.redis.Dedicate()

	go func() {
		defer cancel()

		wait := client.SetPubSubHooks(rueidis.PubSubHooks{
			OnMessage: func(m rueidis.PubSubMessage) {
				logger := e.Logger.With().Str("subscription", m.Channel).Logger()
				logger.Debug().Msg("received message")

				callback(EventMessage{message: m})
			},
		})

		<-wait
	}()

	client.Do(
		context.Background(),
		client.B().Subscribe().Channel(e.topicWithPrefix(event)).Build(),
	)
}

// PSubscribe subscribes to a pattern as returned by the `AllTopics` method on an `Event`. A callback function is to be
// provided which will be called when the event is published, the callback will be passed an `EventMessage` which can be
// used to unmarshall the event.
func (e *Engine) PSubscribe(event Event, callback func(event EventMessage)) {
	e.Logger.Debug().Str("topic", e.allTopicsWithPrefix(event)).Msg("subscribing to topic")

	client, cancel := e.redis.Dedicate()

	go func() {
		defer cancel()

		wait := client.SetPubSubHooks(rueidis.PubSubHooks{
			OnMessage: func(m rueidis.PubSubMessage) {
				logger := e.Logger.With().Str("subscription", m.Channel).Logger()
				logger.Debug().Msg("received message")

				callback(EventMessage{message: m})
			},
		})

		<-wait
	}()

	client.Do(
		context.Background(),
		client.B().Psubscribe().Pattern(e.allTopicsWithPrefix(event)).Build(),
	)
}

// Unsubscribe unsubscribes an event. The event must implement the `Event` interface.
func (e *Engine) Unsubscribe(event Event) {
	e.Logger.Debug().Str("topic", e.topicWithPrefix(event)).Msg("unsubscribing from topic")

	e.redis.Do(
		context.Background(),
		e.redis.B().Unsubscribe().Channel(e.topicWithPrefix(event)).Build(),
	)
}

// PUnsubscribe unsubscribes from a pattern as returned by the `AllTopics` method on an `Event`.
func (e *Engine) PUnsubscribe(event Event) {
	e.Logger.Debug().Str("topic", e.allTopicsWithPrefix(event)).Msg("unsubscribing from topic")

	e.redis.Do(
		context.Background(),
		e.redis.B().Punsubscribe().Pattern(e.allTopicsWithPrefix(event)).Build(),
	)
}

func (e *Engine) topicWithPrefix(event Event) string {
	return fmt.Sprintf("%s:%s", e.instanceId, event.Topic())
}

func (e *Engine) allTopicsWithPrefix(event Event) string {
	return fmt.Sprintf("%s:%s", e.instanceId, event.AllTopics())
}
