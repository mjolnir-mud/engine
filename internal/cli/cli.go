package cli

import "github.com/mjolnir-mud/engine/internal/instance"

type StartCmd struct {
	Service string
}

func (s *StartCmd) Run() {
	instance.StartService(s.Service)
}

var CLI struct {
	Start StartCmd
}
