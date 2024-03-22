package install

import (
	"fmt"
	integrate "github.com/satisfactorymodding/SMEI/cmd/install/wwise"
	"github.com/satisfactorymodding/SMEI/config"
	"github.com/satisfactorymodding/SMEI/lib/cfmt"
	"github.com/satisfactorymodding/SMEI/lib/credentials"
	"github.com/satisfactorymodding/SMEI/lib/elevate"
	"github.com/satisfactorymodding/SMEI/lib/env/project"
	"github.com/satisfactorymodding/SMEI/lib/env/ue"
	"github.com/satisfactorymodding/SMEI/lib/env/vs"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	flags := Cmd.Flags()

	flags.BoolP("local", "l", false, "Install dependencies in the target directory instead of globally")
	flags.StringP("target", "t", "", "Where to install the project")
	flags.BoolP("nonelevated", "e", false, "Choose whether to elevate the process or not. UE installation requires privileges")

	requiredFlags := []string{"target"}
	for _, flag := range requiredFlags {
		err := Cmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatalf("Could not mark flag '%v' as required: %v", flag, err)
		}
	}

	Cmd.AddCommand(integrate.Cmd)
}

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install a modding environment, or components of one",
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

		err := config.Setup()

		err = viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Panicf("Could not bind the CLI flags to the configuration system: %v", err)
		}

		doElevate := !viper.GetBool("nonelevated")
		if doElevate {
			elevate.EnsureElevatedFinal()
		}

		if !config.HasPassword() {
			err = credentials.AskForPassword()
			if err != nil {
				log.Panicf("Could not get a password: %v", err)
			}
		}

		// Collect Wwise credentials in advance of any downloading steps so no interactivity is required mid-install
		wwiseCredentials, err := credentials.GetWwiseCredentials()
		if err != nil {
			log.Panicf("Could not get the Wwise credentials: %v", err)
		}

		local := viper.GetBool("local")
		target := viper.GetString("target")

		cfmt.Sequence.Println("Checking SMEI cached files")
		installerDir := os.TempDir()
		if viper.GetBool(config.PreserveUEInstaller_key) {
			installerDir = filepath.Join(config.ConfigDir, ue.CacheFolder)
		}

		// If lacking github credentials, this will prompt for them. Not needed if the installer files don't need to be downloaded.
		// No further user interaction should be required past this point.
		cfmt.Sequence.Println("Analyzing Unreal Engine install")
		UEInstallDir := viper.GetString(config.UEInstallPath_key)
		fmt.Printf("Expecting UE install dir to be at '%v'\n", UEInstallDir)
		if local {
			UEInstallDir = filepath.Join(target, config.UEFolderName)
		}
		avoidUeReinstall := viper.GetBool(config.UESkipReinstall_key)
		err = ue.Install(UEInstallDir, installerDir, avoidUeReinstall)
		if err != nil {
			log.Panicf("Could not install the Unreal Engine: %v", err)
		}

		cfmt.Sequence.Println("Installing Visual Studio...")
		VSInstallPath := viper.GetString(config.VSInstallPath_key)
		if local {
			VSInstallPath = filepath.Join(target, "VS22")
		}
		avoidVsReinstall := viper.GetBool(config.VSSkipReinstall_key)
		err = vs.Install(VSInstallPath, avoidVsReinstall)
		if err != nil {
			log.Panicf("Could not install Visual Studio: %v", err)
		}

		cfmt.Sequence.Println("Installing modding project...")
		err = project.Install(target, UEInstallDir, *wwiseCredentials)

		if err != nil {
			log.Panicf("Could not install the project: %v", err)
		}
	},
}
