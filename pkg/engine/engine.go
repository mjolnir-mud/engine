package engine

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Init initializes the engine, and execute the CLI.
func Init(name string, plugins []Plugin) {
	setName(name)
	setPlugins(plugins)
	setBaseCommand(&cobra.Command{
		Use:   name,
		Short: "Interact with the Mjolnir MUD engine",
		Long:  `Interact with the Mjolnir MUD engine.`,
	})

	setLogger()
	connectToNats()
	connectToRedis()

	err := loadPlugins()

	if err != nil {
		fmt.Print(fmt.Errorf("error loading plugins: %s", err))
		os.Exit(1)
	}
}
