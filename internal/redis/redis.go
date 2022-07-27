package redis

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Start() *redis.Client {
	// connect to redis
	redisLogger.Info().Msg("Connecting to redis")

	err := viper.BindEnv("redis_url")

	if viper.GetString("env") == "test" {
		viper.SetDefault("redis_url", "redis://localhost:6379/1")
	} else {
		viper.SetDefault("redis_url", "redis://localhost:6379/0")
	}

	if err != nil {
		panic(err)
	}

	u, err := url.Parse(viper.GetString("redis_url"))

	if err != nil {
		redisLogger.Error().Err(err).Msg("error parsing redis url")
		os.Exit(1)
	}

	host := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())

	redisLogger.Info().Msgf("Connecting to redis at %s on %s", host, u.Path)

	i, err := strconv.Atoi(u.Path[1:])

	if err != nil {
		redisLogger.Error().Err(err).Msg("error parsing redis url")
		os.Exit(1)
	}

	client = redis.NewClient(&redis.Options{
		Addr:        host,
		DB:          i,
		PoolSize:    10,
		ReadTimeout: -1,
	})

	// ping redis to ensure it is connected
	_, err = client.Ping(context.Background()).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error connecting to redis")
		os.Exit(1)
	}

	return client
}

func Stop() {
	redisLogger.Info().Msg("Disconnecting from redis")
	err := client.Close()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error disconnecting from redis")
	}
}

func GetClient() *redis.Client {
	return client
}

var client *redis.Client

var redisLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "redis").
	Logger()
