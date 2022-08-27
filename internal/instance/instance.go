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
	Use:   "mjolnir",
	Short: "Mjolnir MUD Engine",
	Long:  "Mjolnir MUD Engine CLI interface",
}

var beforeStartCallbacks = make([]func(), 0)
var afterStartCallbacks = make([]func(), 0)
var beforeStopCallbacks = make([]func(), 0)
var afterStopCallbacks = make([]func(), 0)
var beforeProcessStartCallbacks = make(map[string][]func())
var afterProcessStartCallbacks = make(map[string][]func())

func IsRunning() bool {
	return redis2.Ping() == nil
}

func ExecuteCLI() {
	_ = command.Execute()
}

func RegisterAfterStartCallback(f func()) {
	afterStartCallbacks = append(afterStartCallbacks, f)
}

func RegisterAfterStopCallback(f func()) {
	afterStopCallbacks = append(afterStopCallbacks, f)
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
	setEnv(n)

	logger.Start()
	fmt.Print("Mjolnir MUD Engine\n")
	log = logger.Instance.With().Str("component", "engine").Logger()
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

func setEnv(name string) {
	viper.SetEnvPrefix("MJOLNIR")
	err := viper.BindEnv("env")

	if err != nil {
		panic(err)
	}

	viper.SetDefault("env", "development")

	err = viper.BindEnv("redis_url")

	viper.Set("name", name)

	if err != nil {
		panic(err)
	}
}
