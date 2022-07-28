package command_sets

import (
	movement2 "github.com/mjolnir-mud/engine/plugins/world/internal/commands/movement"
)

type movement struct {
	North     movement2.North     `cmd aliases:"n,N"`
	Northeast movement2.Northeast `cmd aliases:"ne,NE"`
	Northwest movement2.Northwest `cmd aliases:"nw,NW"`
	South     movement2.South     `cmd aliases:"s,S"`
	Southeast movement2.Southeast `cmd aliases:"se,SE"`
	Southwest movement2.Southwest `cmd aliases:"sw,SW"`
	West      movement2.West      `cmd aliases:"w,W"`
	East      movement2.East      `cmd aliases:"e,E"`
	In        movement2.In        `cmd aliases:"in,IN"`
	Out       movement2.Out       `cmd aliases:"out,OUT"`
	Up        movement2.Up        `cmd aliases:"up,UP"`
	Down      movement2.Down      `cmd aliases:"down,DOWN"`
}

func (c *movement) Name() string {
	return "movement"
}

var Movement = &movement{}
