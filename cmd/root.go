package cmd

import (
	configCmd "SMEI/cmd/config"
	"SMEI/cmd/install"
	"SMEI/cmd/test"
	"SMEI/lib/cmdhelp"

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
