package cli

import (
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/server"
	"github.com/spf13/cobra"
)

var CLI = &cobra.Command{
	Use:   "telnet_portal",
	Short: "start the Mjolnir Telnet Portal",
	Long:  "start the Mjolnir Telnet Portal",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New()

		s.Start()
	},
}
