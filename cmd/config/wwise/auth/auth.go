package auth

import (
	"SMEI/lib/cmdhelp"
	"SMEI/lib/colors"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "UNIMPLEMENTED Configure Wwise authentication",
	Run: func(cmd *cobra.Command, args []string) {
		colors.Sequence.Println("This command is not yet implemented")
		cmdhelp.PrintHelp(cmd)
	},
}
