package integrate

import (
	"SMEI/config"
	"SMEI/lib/colors"
	"SMEI/lib/credentials"
	"SMEI/lib/env/project"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	flags := Cmd.Flags()

	flags.StringP("target", "t", "", "Path to existing project folder")

	requiredFlags := []string{"target"}
	for _, flag := range requiredFlags {
		err := Cmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatalf("Could not mark flag '%v' as required: %v", flag, err)
		}
	}
}

var Cmd = &cobra.Command{
	Use:   "integrate",
	Short: "Integrate wwise into an existing project",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			v := recover()
			if v != nil {
				fmt.Println(v)
			}
			fmt.Println("Use ctrl+C to close this window")
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt)
			<-c
		}()

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Panicf("Could not bind the CLI flags to the configuration system: %v", err)
		}

		target := viper.GetString("target")

		err = config.Setup()

		if !config.HasPassword() {
			err = credentials.AskForPassword()
			if err != nil {
				log.Panicf("Could not get a password: %v", err)
			}
		}

		wwiseCredentials, err := credentials.GetWwiseCredentials()
		if err != nil {
			log.Panicf("Could not get the Wwise credentials: %v", err)
		}

		colors.Sequence.Printf("Integrating Wwise into project '%s'...\n", target)
		err = project.InstallWWise(target, *wwiseCredentials)

		if err != nil {
			log.Panicf("Could not integrate wwise the project: %v", err)
		}
	},
}
