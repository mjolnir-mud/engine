package redis

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var Client *redis.Client
var log zerolog.Logger

type RedisLogProxy struct {
	logger zerolog.Logger
}

func (l *RedisLogProxy) Printf(_ context.Context, format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func CreateAndStart() {
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

	Client = redis.NewClient(&redis.Options{
		Addr:        host,
		DB:          i,
		PoolSize:    10,
		ReadTimeout: -1,
	})

}

func Stop() {
	_ = Client.Close()
}
