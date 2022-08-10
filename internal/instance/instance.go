package instance

import (
	"fmt"

	"github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	redis2 "github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func Initialize(n string) {
	name = n
	setEnv()
	logger.Start()
	fmt.Print("Mjolnir MUD Engine\n")
	log = logger.Instance
	log.Info().Msgf("initializing engine for game %s", name)
	redis2.CreateAndStart()
	Redis = redis2.Client
}

func Start() {
	plugin_registry.Start()
}

func Stop() {
	log.Info().Msg("stopping engine")
	plugin_registry.Stop()
	_ = Redis.Close()
}

var name string
var Redis *redis.Client
var log zerolog.Logger

func setEnv() {
	viper.SetEnvPrefix("MJOLNIR")
	err := viper.BindEnv("env")

	if err != nil {
		panic(err)
	}

	viper.SetDefault("env", "development")

	err = viper.BindEnv("redis_url")

	if err != nil {
		panic(err)
	}
}
