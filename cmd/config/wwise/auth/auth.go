package auth

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure Wwise authentication",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
