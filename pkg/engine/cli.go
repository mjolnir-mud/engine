package engine

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AddCommand can be called by plugins to add a command to the base command.
func AddCommand(command *cobra.Command) {
	state.baseCommand.AddCommand(command)
}

func ExecCommand() {
	err := state.baseCommand.Execute()

	if err != nil {
		fmt.Print(fmt.Errorf("error executing command: %s", err))
		os.Exit(1)
	}
}

func setBaseCommand(command *cobra.Command) {
	state.baseCommand = command
}
