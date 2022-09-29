/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package command_registry

import (
	"github.com/mjolnir-mud/engine/plugins/sessions/systems/session"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/mjolnir-mud/engine/plugins/command_parser/pkg/command_set"
)

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

func ParseCommand(sets []string, sessId string, input string) {
	for _, set := range sets {
		// get the command set if mil write  an error to the connection
		k, ok := registry.commandSets[set]

		if !ok {
			err := session.SendLine(sessId, "I'm confused.")

			if err != nil {
				panic(err)
			}

			return
		}

		// parse the command
		ctx, err := k.Parse(strings.Split(input, " "))

		if err != nil {
			continue
		} else {
			ctx.Bind(sessId)

			err = ctx.Run(sessId)

			if err != nil {
				err := session.SendLine(sessId, err.Error())

				if err != nil {
					panic(err)
				}

				return
			}

			return
		}
	}
	err := session.SendLine(sessId, "I'm not sure what you mean.")

	if err != nil {
		panic(err)
	}
}
