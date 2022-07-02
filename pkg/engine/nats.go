package engine

import (
	"os"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

func connectToNats() {
	log.Info().Msg("Connecting to nats")
	n, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to nats")
		os.Exit(1)
	}

	c, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)

	if err != nil {
		log.Error().Err(err).Msg("Error connecting to nats")
		os.Exit(1)
	}

	state.natsConn = n
	state.nats = c

	log.Debug().Msg("Connected to nats")
}

// PublishEvent publishes an event to the NATs event bus.
func PublishEvent(topic string, event interface{}) error {
	log.Debug().Msgf("Publishing event %s", topic)
	err := state.nats.Publish(topic, event)

	if err != nil {
		log.Error().Err(err).Msg("Error publishing event")
		return err
	}

	return nil
}

// SubscribeToEvent subscribes to an event on the NATS event bus.
func SubscribeToEvent(topic string, handler nats.Handler) (*nats.Subscription, error) {
	logger := log.With().Str("topic", topic).Logger()
	logger.Debug().Msg("Subscribing to event")

	sub, err := state.nats.Subscribe(topic, handler)

	if err != nil {
		logger.Error().Err(err).Msg("Error subscribing to event")
		return nil, err
	}

	logger.Trace().Msg("Subscribed to event")

	return sub, nil
}
