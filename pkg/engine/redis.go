package engine

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

func connectToRedis() {
	log.Info().Msg("Connecting to Redis")
	state.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// PublishEvent publishes an event to the Redis event bus after serializing the event as JSON.
func PublishEvent(context context.Context, topic string, event interface{}) error {
	j, err := json.Marshal(event)

	if err != nil {
		log.Error().Err(err).Msg("Error marshalling event")
		return err
	}

	log.Debug().Msgf("Publishing event %s", topic)
	state.rdb.Publish(context, topic, j)

	return nil
}

// SubscribeToEvent subscribes to an event on the Redis event bus.
func SubscribeToEvent(context context.Context, topic string, handler func(context.Context, *EventMessage)) {
	log.Debug().Msgf("Subscribing to event %s", topic)
	ch := state.rdb.Subscribe(context, topic).Channel()

	go func() {
		for msg := range ch {
			b := msg.Payload

			event := NewEventMessage(topic, b)

			handler(context, event)
		}
	}()
}
