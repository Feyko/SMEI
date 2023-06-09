package wwise

import (
	"SMEI/lib/cmdhelp"
	"SMEI/lib/colors"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "wwise",
	Short: "UNIMPLEMENTED Configure Wwise",
	Run: func(cmd *cobra.Command, args []string) {
		colors.Sequence.Println("This command is not yet implemented")
		cmdhelp.PrintHelp(cmd)
	},
}
