package cmdhelp

import (
	"github.com/spf13/cobra"
	"log"
)

func PrintHelp(cmd *cobra.Command) {
	err := cmd.Help()
	if err != nil {
		log.Fatalf("Could not print the help: %v", err)
	}
}
