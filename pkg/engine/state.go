package engine

import (
	"github.com/spf13/cobra"
)

// State represents the state of the game.
type State struct {
	name        string
	baseCommand *cobra.Command
}

func Name() string {
	return state.name
}

func setName(name string) {
	state.name = name
}

var state = &State{}
