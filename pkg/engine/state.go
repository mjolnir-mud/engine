package engine

import (
	"github.com/go-redis/redis/v9"
	"github.com/spf13/cobra"
)

// State represents the state of the game.
type State struct {
	name        string
	plugins     []Plugin
	baseCommand *cobra.Command
	rdb         *redis.Client
}

func Name() string {
	return state.name
}

func setName(name string) {
	state.name = name
}

func setPlugins(plugins []Plugin) {
	state.plugins = plugins
}

var state = &State{}
