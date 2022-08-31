package instance

import (
	"fmt"
	"github.com/mjolnir-mud/engine/internal/plugin_registry"
	redis2 "github.com/mjolnir-mud/engine/internal/redis"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

var beforeStartCallbacks = make([]func(), 0)
var afterStartCallbacks = make([]func(), 0)
var beforeStopCallbacks = make([]func(), 0)
var afterStopCallbacks = make([]func(), 0)
var onServiceStartCallbacks = make(map[string][]func())
var onServiceStopCallbacks = make(map[string][]func())
var onEnvStartCallbacks = make(map[string][]func())
var onEnvStopCallbacks = make(map[string][]func())
var environment string
var gameName string

var Running chan bool

func SetEnv(n string) {
	environment = n
}

func GetEnv() string {
	return environment
}

func GetGameName() string {
	return gameName
}

func IsRunning() bool {
	return redis2.Ping() == nil
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

func RegisterOnServiceStartCallback(service string, f func()) {
	onServiceStartCallbacks[service] = append(onServiceStartCallbacks[service], f)
}

func RegisterOnServiceStopCallback(service string, f func()) {
	onServiceStopCallbacks[service] = append(onServiceStopCallbacks[service], f)
}

func RegisterOnEnvStartCallback(env string, f func()) {
	onEnvStartCallbacks[env] = append(onEnvStartCallbacks[env], f)
}

func RegisterOnEnvStopCallback(env string, f func()) {
	onEnvStopCallbacks[env] = append(onEnvStopCallbacks[env], f)
}

func Start(n string) {
	gameName = n

	logger.Start()
	fmt.Print("Mjolnir MUD Engine\n")
	log = logger.Instance.With().Str("component", "engine").Logger()
	log.Info().Msgf("running beforeStartCallbacks")
	redis2.Start()

	for _, f := range beforeStartCallbacks {
		f()
	}

	env := GetEnv()

	for _, f := range onEnvStartCallbacks[env] {
		f()
	}

	plugin_registry.Start()
	log.Info().Msgf("running afterStartCallbacks")
	for _, f := range afterStartCallbacks {
		f()
	}
}

func StopService(service string) {
	log.Info().Str("service", service).Msg("stopping service")
	for _, f := range onServiceStopCallbacks[service] {
		f()
	}
}

func StartService(service string) {
	log.Info().Str("service", service).Msg("starting service")

	Running = make(chan bool)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		StopService(service)
		Stop()
		Running <- true
	}()

	for _, f := range onServiceStartCallbacks[service] {
		f()
	}
}

func Stop() {
	log.Info().Msg("stopping engine")
	for _, f := range beforeStopCallbacks {
		f()
	}

	env := GetEnv()

	for _, f := range onEnvStopCallbacks[env] {
		f()
	}

	plugin_registry.Stop()
	redis2.Stop()

	for _, f := range afterStopCallbacks {
		f()
	}
}

var log zerolog.Logger
