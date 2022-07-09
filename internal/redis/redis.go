package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Start() {
	// connect to redis
	redisLogger.Info().Msg("Connecting to redis")

	err := viper.BindEnv("redis_url")
	viper.SetDefault("redis_url", "redis://localhost:6379/0")

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

	redisClient = redis.NewClient(&redis.Options{
		Addr: host,
		DB:   i,
	})

	// ping redis to ensure it is connected
	_, err = redisClient.Ping(context.Background()).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error connecting to redis")
		os.Exit(1)
	}
}

func Stop() {
	redisLogger.Info().Msg("Disconnecting from redis")
	err := redisClient.Close()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error disconnecting from redis")
	}
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	redisLogger.Debug().Msgf("Setting key %s", key)
	// marshall value to json
	jsonValue, err := json.Marshal(value)

	if err != nil {
		redisLogger.Error().Err(err).Msg("error marshalling value")
		return err
	}

	err = redisClient.Set(context.Background(), key, jsonValue, expiration).Err()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error setting key")
		return err
	}

	return nil
}

func RedisGet(key string, value interface{}) error {
	redisLogger.Debug().Msgf("Getting key %s", key)
	jsonValue, err := redisClient.Get(context.Background(), key).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error getting key")
		return err
	}

	err = json.Unmarshal([]byte(jsonValue), value)

	if err != nil {
		redisLogger.Error().Err(err).Msg("error unmarshalling value")
		return err
	}

	return nil
}

func RedisKeyExists(key string) bool {
	redisLogger.Debug().Msg("Checking if key exists")
	exists, err := redisClient.Exists(context.Background(), key).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error checking if key exists")
		os.Exit(1)
	}

	if exists > 0 {
		return true
	} else {
		return false
	}
}

var redisClient *redis.Client
var redisLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "redis").
	Logger()
