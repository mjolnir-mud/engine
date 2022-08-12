package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var log zerolog.Logger

type RedisLogProxy struct {
	logger zerolog.Logger
}

var client *redis.Client

func (l *RedisLogProxy) Printf(_ context.Context, format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func Ping() error {
	return client.Ping(context.Background()).Err()
}

func Publish(e event.Event, args ...interface{}) error {
	p := e.Payload(args...)

	payloadBytes, err := json.Marshal(p)

	topic := e.Topic(args...)

	if err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("error marshalling event")
		return err
	}

	log.Debug().Str("topic", topic).Msgf("publishing event: %d", len(payloadBytes))
	err = client.Publish(context.Background(), topic, string(payloadBytes)).Err()

	if err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("error publishing event")
		return err
	}

	return nil
}

func Start() {
	if viper.GetString("env") == "test" {
		viper.SetDefault("redis_url", "redis://localhost:6379/1")
	} else {
		viper.SetDefault("redis_url", "redis://localhost:6379/0")
	}

	u, err := url.Parse(viper.GetString("redis_url"))

	log = logger.Instance.
		With().
		Str("service", "redis").
		Logger()

	if err != nil {
		log.Fatal().Err(err).Msg("could not parse redis url")
		os.Exit(1)
	}

	host := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())

	log.Info().Msgf("connecting to redis at %s on %s", host, u.Path)

	i, err := strconv.Atoi(u.Path[1:])

	if err != nil {
		log.Fatal().Err(err).Msg("could not parse redis url")
		os.Exit(1)
	}

	log = log.With().
		Str("host", host).
		Int("db", i).
		Logger()

	redis.SetLogger(&RedisLogProxy{log})

	client = redis.NewClient(&redis.Options{
		Addr:        host,
		DB:          i,
		PoolSize:    10,
		ReadTimeout: -1,
	})

}

func Stop() {
	_ = client.Close()
}
