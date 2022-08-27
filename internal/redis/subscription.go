package redis

import (
	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
)

type Subscription struct {
	pubsub   *redis.PubSub
	stop     chan bool
	callback func(payload event.EventPayload)
	logger   zerolog.Logger
	event    event.Event
}

func NewSubscription(e event.Event, callback func(payload event.EventPayload)) *Subscription {
	return createSubscription(Subscribe(e.Topic()), e, callback)

}

func NewPatternSubscription(e event.Event, callback func(payload event.EventPayload)) *Subscription {
	return createSubscription(PSubscribe(e.Topic()), e, callback)
}

func createSubscription(pubsub *redis.PubSub, e event.Event, callback func(payload event.EventPayload)) *Subscription {
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

				newEvent := event.EventPayload{
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
