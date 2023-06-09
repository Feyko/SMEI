package wwise

import (
	"SMEI/lib/cfmt"
	"SMEI/lib/cmdhelp"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "wwise",
	Short: "UNIMPLEMENTED Configure Wwise",
	Run: func(cmd *cobra.Command, args []string) {
		cfmt.Sequence.Println("This command is not yet implemented")
		cmdhelp.PrintHelp(cmd)
	},
}
