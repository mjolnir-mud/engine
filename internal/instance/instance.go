package instance

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	redis2 "github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var command = &cobra.Command{
	Use: "mjolnir",
}
var beforeStartCallbacks = make([]func(), 0)
var afterStartCallbacks = make([]func(), 0)
var beforeStopCallbacks = make([]func(), 0)
var afterStopCallbacks = make([]func(), 0)

func IsRunning() bool {
	return redis2.Ping() == nil
}

func RegisterAfterStartCallback(f func()) {
	afterStartCallbacks = append(afterStartCallbacks, f)
}

func RegisterAfterStopCallback(f func()) {
	afterStopCallbacks = append(beforeStopCallbacks, f)
}

func RegisterBeforeStopCallback(f func()) {
	beforeStopCallbacks = append(beforeStopCallbacks, f)
}

func RegisterBeforeStartCallback(f func()) {
	beforeStartCallbacks = append(beforeStartCallbacks, f)
}

func RegisterCLICommand(c *cobra.Command) {
	command.AddCommand(c)
}

func Start(n string) {
	name = n
	setEnv()
	logger.Start()
	fmt.Print("Mjolnir MUD Engine\n")
	log = logger.Instance
	log.Info().Msgf("running beforeStartCallbacks")
	redis2.Start()
	for _, f := range beforeStartCallbacks {
		f()
	}

	plugin_registry.Start()
	log.Info().Msgf("running afterStartCallbacks")
	for _, f := range afterStartCallbacks {
		f()
	}
}

func Stop() {
	log.Info().Msg("stopping engine")
	for _, f := range beforeStopCallbacks {
		f()
	}
	plugin_registry.Stop()
	redis2.Stop()

	for _, f := range afterStopCallbacks {
		f()
	}
}

var name string
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
