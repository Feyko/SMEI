package auth

import (
	"github.com/satisfactorymodding/SMEI/lib/cfmt"
	"github.com/satisfactorymodding/SMEI/lib/cmdhelp"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "UNIMPLEMENTED Configure Wwise authentication",
	Run: func(cmd *cobra.Command, args []string) {
		cfmt.Sequence.Println("This command is not yet implemented")
		cmdhelp.PrintHelp(cmd)
	},
}
