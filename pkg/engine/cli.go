package engine

import "github.com/spf13/cobra"

// AddCommand can be called by plugins to add a command to the base command.
func AddCommand(command *cobra.Command) {
	state.baseCommand.AddCommand(command)
}

func setCommand(command *cobra.Command) {
	state.baseCommand = command
}
