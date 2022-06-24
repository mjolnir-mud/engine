package engine

import (
	"fmt"
	"os"
)

// Init initializes the engine, and execute the CLI.
func Init(plugins []Plugin) {
	state := newState(plugins)

	err := loadPlugins(state)

	if err != nil {
		fmt.Print(fmt.Errorf("error loading plugins: %s", err))
		os.Exit(1)
	}
}
