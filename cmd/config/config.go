package config

import (
	"SMEI/cmd/config/wwise"
	"SMEI/cmdhelp"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Configure SMEI",
	Run: func(cmd *cobra.Command, args []string) {
		cmdhelp.PrintHelp(cmd)
	},
}

func init() {
	Cmd.AddCommand(wwise.Cmd)
}
