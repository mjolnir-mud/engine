package telnet_portal

import (
	"os"

	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/server"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/pkg/config"
)

type telnetPortal struct {
}

func (p telnetPortal) Name() string {
	return "telnet_portal"
}

func (p telnetPortal) Registered() error {
	engine.RegisterAfterServiceStartCallback("telnet_portal", func() {
		env := os.Getenv("MJOLNIR_ENV")

		if env == "" {
			env = "default"
		}

		cfgCB, ok := config.Configurations[env]

		if !ok {
			panic("no configuration found for environment: " + env)
		}

		cfg := cfgCB(&config.Configuration{})

		if cfg == nil {
			panic("configuration callback returned nil")
		}

		server.Start(cfg)
	})

	return nil
}

func ConfigureForEnv(env string, cb func(configuration *config.Configuration) *config.Configuration) {
	defaultConfig := config.Configurations["default"](&config.Configuration{})
	config.Configurations[env] = func(configuration *config.Configuration) *config.Configuration {
		return cb(defaultConfig)
	}
}

var Plugin = telnetPortal{}
