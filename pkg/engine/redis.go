package engine

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

func connectToRedis() {
	// connect to redis
	log.Info().Msg("Connecting to redis")
	state.redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// ping redis
	_, err := state.redis.Ping(context.Background()).Result()

	if err != nil {
		log.Error().Err(err).Msg("error connecting to redis")
		os.Exit(1)
	}
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	log.Debug().Msgf("Setting key %s", key)
	// marshall value to json
	jsonValue, err := json.Marshal(value)

	if err != nil {
		log.Error().Err(err).Msg("error marshalling value")
		return err
	}

	err = state.redis.Set(context.Background(), key, jsonValue, expiration).Err()

	if err != nil {
		log.Error().Err(err).Msg("error setting key")
		return err
	}

	return nil
}

func RedisGet(key string, value interface{}) error {
	log.Debug().Msgf("Getting key %s", key)
	jsonValue, err := state.redis.Get(context.Background(), key).Result()

	if err != nil {
		log.Error().Err(err).Msg("error getting key")
		return err
	}

	err = json.Unmarshal([]byte(jsonValue), value)

	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling value")
		return err
	}

	return nil
}

func RedisKeyExists(key string) bool {
	log.Debug().Msg("Checking if key exists")
	exists, err := state.redis.Exists(context.Background(), key).Result()

	if err != nil {
		log.Error().Err(err).Msg("error checking if key exists")
		os.Exit(1)
	}

	if exists > 0 {
		return true
	} else {
		return false
	}
}
