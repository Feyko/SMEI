package cmd

import (
	"SMEI/config"
	"SMEI/lib/project"
	"SMEI/lib/secret"
	"SMEI/lib/ue"
	"SMEI/lib/vs"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "smei",
	Short: "Assists in setting up a modding environment for Satisfactory",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		config.Setup()

		err := viper.BindPFlags(flags)
		if err != nil {
			log.Fatalf("Could not bind the CLI flags to the configuration system: %v", err)
		}

		local := viper.GetBool("local")
		target := viper.GetString("target")

		UEInstallDir := viper.GetString(config.UEInstallPath_key)
		if local {
			UEInstallDir = filepath.Join(target, ue.FolderName)
		}
		installerDir := os.TempDir()
		if viper.GetBool(config.PreserveUEInstaller_key) {
			installerDir = filepath.Join(config.ConfigDir, "UE-Installer")
		}

		fmt.Println("Installing the Unreal Engine")
		err = ue.Install(UEInstallDir, installerDir)
		if err != nil {
			log.Fatalf("Could not install the Unreal Engine: %v", err)
		}

		fmt.Println("Installing Visual Studio...")
		VSInstallPath := viper.GetString(config.VSInstallPath_key)
		if local {
			VSInstallPath = filepath.Join(config.ConfigDir, "Visual Studio 2022 - Community")
		}
		err = vs.Install(VSInstallPath)
		if err != nil {
			log.Fatalf("Could not install Visual Studio: %v", err)
		}

		fmt.Println("Installing modding project...")
		wwiseEmail := secret.String(viper.GetString(config.WwiseEmail_key))
		if wwiseEmail == "" {
			log.Fatalf("Wwise email not provided")
		}
		wwisePassword := secret.String(viper.GetString(config.WwisePassword_key))
		if wwisePassword == "" {
			log.Fatalf("Wwise password not provided")
		}

		err = project.Install(target, UEInstallDir, project.WwiseAuth{
			Email:    wwiseEmail,
			Password: wwisePassword,
		})

		if err != nil {
			log.Fatalf("Could not install the project: %v", err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	flags := rootCmd.Flags()

	flags.BoolP("local", "l", false, "Install dependencies in the target directory instead of globally")
	flags.StringP("target", "t", "", "Where to install the project")

	requiredFlags := []string{"target"}
	for _, flag := range requiredFlags {
		err := rootCmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatalf("Could not mark flag '%v' as required: %v", flag, err)
		}
	}
}
