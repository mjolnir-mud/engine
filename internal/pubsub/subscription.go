package pubsub

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	redis2 "github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
)

type Subscription struct {
	pubsub   *redis.PubSub
	stop     chan bool
	callback func(payload interface{})
	logger   zerolog.Logger
	event    func() interface{}
}

func Subscribe(topic string, event func() interface{}, callback func(payload interface{})) *Subscription {
	return create(redis2.Client.Subscribe(context.Background(), topic), event, callback)
}

func PSubscribe(topic string, event func() interface{}, callback func(payload interface{})) *Subscription {
	return create(redis2.Client.PSubscribe(context.Background(), topic), event, callback)
}

func create(pubsub *redis.PubSub, event func() interface{}, callback func(paylaod interface{})) *Subscription {
	s := &Subscription{
		pubsub:   pubsub,
		stop:     make(chan bool),
		callback: callback,
		event:    event,
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

				newEvent := s.event()

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
