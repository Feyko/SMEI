package wwise

import (
	"SMEI/cmdhelp"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "wwise",
	Short: "Configure Wwise",
	Run: func(cmd *cobra.Command, args []string) {
		cmdhelp.PrintHelp(cmd)
	},
}
