package config

import (
	"SMEI/cmd/config/wwise"
	"SMEI/lib/cfmt"
	"SMEI/lib/cmdhelp"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "UNIMPLEMENTED Configure SMEI",
	Run: func(cmd *cobra.Command, args []string) {
		cfmt.Sequence.Println("This command is not yet implemented. Config file is stored at '%APPDATA%\\SMEI\\'")
		cmdhelp.PrintHelp(cmd)
	},
}

func init() {
	Cmd.AddCommand(wwise.Cmd)
}
