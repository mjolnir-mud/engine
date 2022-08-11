package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
)

type Subscription struct {
	pubsub   *redis.PubSub
	stop     chan bool
	callback func(payload interface{})
	logger   zerolog.Logger
	event    event.Event
}

func Subscribe(e event.Event, args ...interface{}) *Subscription {
	callback, ok := args[len(args)-1].(func(payload interface{}))
	if !ok {
		panic("callback is not a function")
	}
	// remove the last argument as the callback
	args = args[:len(args)-1]

	return createSubscription(client.Subscribe(context.Background(), e.Topic(args...)), e, callback)

}

func PSubscribe(e event.Event, args ...interface{}) *Subscription {
	callback, ok := args[len(args)-1].(func(payload interface{}))
	if !ok {
		panic("callback is not a function")
	}
	// remove the last argument as the callback
	args = args[:len(args)-1]

	return createSubscription(client.PSubscribe(context.Background(), e.Topic(args...)), e, callback)
}

func createSubscription(pubsub *redis.PubSub, e event.Event, callback func(payload interface{})) *Subscription {
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
		for {
			select {
			case <-s.stop:
				_ = s.pubsub.Close()
				return
			case msg := <-s.pubsub.Channel():
				payloadBytes := []byte(msg.Payload)
				length := len(payloadBytes)

				s.logger.Debug().Msgf("received message: %d", length)

				newEvent := e.Payload()

				err := json.Unmarshal([]byte(msg.Payload), newEvent)

				if err != nil {
					s.logger.Error().Err(err).Msg("error unmarshalling message")
					continue
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
