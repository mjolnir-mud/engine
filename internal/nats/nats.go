package nats

import (
	"os"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Start() {
	err := viper.BindEnv("nats_url")

	if err != nil {
		panic(err)
	}

	viper.SetDefault("nats_url", nats.DefaultURL)

	natsLogger.Info().Msgf("connecting to nats on %s", viper.GetString("nats_url"))
	n, err := nats.Connect(viper.GetString("nats_url"))

	if err != nil {
		natsLogger.Error().Err(err).Msg("error connecting to nats")
		os.Exit(1)
	}

	c, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)

	if err != nil {
		natsLogger.Error().Err(err).Msg("error connecting to nats")
		os.Exit(1)
	}

	natsConn = c

	natsLogger.Debug().Msg("connected to nats")
}

func Stop() {
	natsLogger.Debug().Msg("disconnecting from nats")
	natsConn.Close()
}

// PublishEvent publishes an event to the NATs event bus.
func PublishEvent(topic string, event interface{}) error {
	natsLogger.Debug().Msgf("publishing event %s", topic)
	err := natsConn.Publish(topic, event)

	if err != nil {
		natsLogger.Error().Err(err).Msg("error publishing event")
		return err
	}

	return nil
}

// SubscribeToEvent subscribes to an event on the NATS event bus.
func SubscribeToEvent(topic string, handler nats.Handler) (*nats.Subscription, error) {
	logger := natsLogger.With().Str("topic", topic).Logger()
	logger.Debug().Msg("subscribing to event")

	sub, err := natsConn.Subscribe(topic, handler)

	if err != nil {
		logger.Error().Err(err).Msg("error subscribing to event")
		return nil, err
	}

	logger.Trace().Msg("subscribed to event")

	return sub, nil
}

var natsConn *nats.EncodedConn
var natsLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "engine").
	Logger()
