package cmd

import (
	configCmd "github.com/satisfactorymodding/SMEI/cmd/config"
	"github.com/satisfactorymodding/SMEI/cmd/install"
	"github.com/satisfactorymodding/SMEI/cmd/test"
	"github.com/satisfactorymodding/SMEI/lib/cmdhelp"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "smei",
	Short: "Assists in setting up a modding environment for Satisfactory",
	Run: func(cmd *cobra.Command, args []string) {
		cmdhelp.PrintHelp(cmd)
	},
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.AddCommand(configCmd.Cmd, install.Cmd)
	if test.Cmd != nil {
		RootCmd.AddCommand(test.Cmd)
	}
}
