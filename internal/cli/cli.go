package cli

import (
	"github.com/mjolnir-mud/engine/internal/instance"
)

type StartCmd struct {
	Service string `required:""`
}

func (s *StartCmd) Run() error {
	instance.StartService(s.Service)

	return nil
}

var CLI struct {
	Start StartCmd
}
