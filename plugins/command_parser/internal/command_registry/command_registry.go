package command_registry

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mjolnir-mud/engine/plugins/command_parser/pkg/command_set"
)

type Session interface {
	WriteToConnection(string)
}

type commandRegistry struct {
	commandSets map[string]*kong.Kong
}

var registry = &commandRegistry{
	commandSets: make(map[string]*kong.Kong),
}

func RegisterCommandSet(set command_set.CommandSet) {
	k, err := kong.New(set)

	if err != nil {
		panic(err)
	}

	registry.commandSets[set.Name()] = k
}

func ParseCommand(sets []string, sess Session, input string) {
	for _, set := range sets {
		// get the command set if mil write  an error to the connection
		k, ok := registry.commandSets[set]

		if !ok {
			sess.WriteToConnection("I'm confused.")
			return
		}

		// parse the command
		ctx, err := k.Parse(strings.Split(input, " "))

		if err != nil {
			continue
		} else {
			ctx.Bind(sess)

			err = ctx.Run(sess)

			if err != nil {
				sess.WriteToConnection(err.Error())
				return
			}

			return
		}
	}
	sess.WriteToConnection("I'm not sure what you mean.")
}
