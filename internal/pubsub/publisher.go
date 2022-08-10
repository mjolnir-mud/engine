package pubsub

import (
	"context"
	"encoding/json"

	"github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/logger"
)

var log = logger.Instance.With().Str("service", "pubsub").Logger()

func Publish(topic string, event interface{}) error {
	payload, err := json.Marshal(event)

	if err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("error marshalling event")
		return err
	}

	log.Debug().Str("topic", topic).Msgf("publishing event: %d", len(payload))
	err = redis.Client.Publish(context.Background(), topic, string(payload)).Err()

	if err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("error publishing event")
		return err
	}

	return nil
}
