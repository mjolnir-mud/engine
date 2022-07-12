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

func Set(key string, value interface{}, expiration time.Duration) error {
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

func Get(key string, value interface{}) error {
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

func Delete(key string) error {
	redisLogger.Debug().Msgf("Deleting key %s", key)
	err := redisClient.Del(context.Background(), key).Err()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error deleting key")
		return err
	}

	return nil
}

func KeyExists(key string) (bool, error) {
	redisLogger.Debug().Msg("Checking if key exists")
	exists, err := redisClient.Exists(context.Background(), key).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error checking if key exists")
		return false, err
	}

	if exists > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func Count(key string) (int64, error) {
	redisLogger.Trace().Msgf("Counting key: %s", key)
	count, err := redisClient.LLen(context.Background(), key).Result()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error counting key")
		os.Exit(1)
	}

	return count, nil
}

func DeleteAll(key string) error {
	redisLogger.Trace().Msgf("Deleting all keys for %s", key)
	err := redisClient.Del(context.Background(), key).Err()

	if err != nil {
		redisLogger.Error().Err(err).Msg("error deleting all keys")
		return err
	}

	return nil
}

var redisClient *redis.Client
var redisLogger = log.
	With().
	Str("plugin", "engine").
	Str("service", "redis").
	Logger()
