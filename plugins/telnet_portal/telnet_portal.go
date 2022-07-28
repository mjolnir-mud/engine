package telnet_portal

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/cli"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/logger"
)

type telnetPortal struct {
}

var log = logger.Logger

func (p telnetPortal) Name() string {
	return "Engine"
}

func (p telnetPortal) Start() error {

	log.Info().Msg("initializing")
	engine.AddCLICommand(cli.CLI)

	return nil
}

var Plugin = &telnetPortal{}
