package engine

import (
	"github.com/spf13/cobra"
)

// Init initializes the engine, and execute the CLI.
func Init(name string, plugs []Plugin) {
	setName(name)
	setBaseCommand(&cobra.Command{
		Use:   name,
		Short: "Interact with the Mjolnir MUD engine",
		Long:  `Interact with the Mjolnir MUD engine.`,
	})

	setLogger()
	connectToNats()
	connectToRedis()

}

func RegisterPlugin(plugin Plugin) {
	plugins.Register(plugin)
}
