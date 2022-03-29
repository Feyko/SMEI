package cmd

import (
	"SMEI/internal/ghauth"
	"SMEI/internal/project"
	"SMEI/internal/ue"
	"SMEI/internal/vs"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "smei",
	Short: "Assists in setting up a modding environment for Satisfactory",
	Run: func(cmd *cobra.Command, args []string) {
		local, err := cmd.Flags().GetBool("local")
		if err != nil {
			log.Fatalf("Could not get the 'local' flag: %v", err)
		}

		target, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalf("Could not get the 'local' flag: %v", err)
		}

		fmt.Println("Authenticating with GitHub")
		token, err := ghauth.GetToken()
		if err != nil {
			log.Fatalf("Could not authenticate with GitHub: %v", err)
		}
		//fmt.Println(token)

		fmt.Println("Installing the Unreal Engine")
		err = ue.Install(token, target, local)
		if err != nil {
			log.Fatalf("Could not install the Unreal Engine: %v", err)
		}

		fmt.Println("Installing Visual Studio...")
		err = vs.Install(local, target)
		if err != nil {
			log.Fatalf("Could not install Visual Studio: %v", err)
		}

		fmt.Println("Installing modding project...")
		err = project.Install(target, local)

		if err != nil {
			log.Fatalf("Could not install the project: %v", err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP("local", "l", false,
		"Install dependencies in the target directory instead of globally")
	rootCmd.Flags().StringP("target", "t", "",
		"Where to install the project")

	requiredFlags := []string{"target"}
	for _, flag := range requiredFlags {
		err := rootCmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatalf("Could not mark flag '%v' as required: %v", flag, err)
		}
	}
}
