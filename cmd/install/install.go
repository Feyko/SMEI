package install

import (
	"SMEI/config"
	"SMEI/lib/elevate"
	"SMEI/lib/env/project"
	"SMEI/lib/env/ue"
	"SMEI/lib/env/vs"
	"SMEI/lib/secret"
	"fmt"
	"github.com/fatih/color"
	"github.com/mircearoata/wwise-cli/lib/wwise/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
}

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install a modding environment",
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
			err = askForPassword()
			if err != nil {
				log.Panicf("Could not get a password: %v", err)
			}
		}

		if !viper.IsSet(config.WwiseEmail_key) {
			err = askForWwiseAuth()
			if err != nil {
				log.Panicf("Could not log in with Wwise: %v", err)
			}
		}

		wwiseEmail, err := config.GetSecretString(config.WwiseEmail_key)
		if err != nil {
			log.Panicf("Could not get the Wwise email: %v", err)
		}

		wwisePassword, err := config.GetSecretString(config.WwisePassword_key)
		if err != nil {
			log.Panicf("Could not get the Wwise password: %v", err)
		}

		local := viper.GetBool("local")
		target := viper.GetString("target")

		UEInstallDir := viper.GetString(config.UEInstallPath_key)
		if local {
			UEInstallDir = filepath.Join(target, config.UEFolderName)
		}
		installerDir := os.TempDir()
		if viper.GetBool(config.PreserveUEInstaller_key) {
			installerDir = filepath.Join(config.ConfigDir, ue.CacheFolder)
		}

		fmt.Println("Installing the Unreal Engine")
		err = ue.Install(UEInstallDir, installerDir)
		if err != nil {
			log.Panicf("Could not install the Unreal Engine: %v", err)
		}

		fmt.Println("Installing Visual Studio...")
		VSInstallPath := viper.GetString(config.VSInstallPath_key)
		if local {
			VSInstallPath = filepath.Join(target, "VS22")
		}
		err = vs.Install(VSInstallPath)
		if err != nil {
			log.Panicf("Could not install Visual Studio: %v", err)
		}

		fmt.Println("Installing modding project...")

		err = project.Install(target, UEInstallDir, project.WwiseAuth{
			Email:    wwiseEmail,
			Password: wwisePassword,
		})

		if err != nil {
			log.Panicf("Could not install the project: %v", err)
		}
	},
}

func askForPassword() error {
	if config.HasLoggedInBefore() {
		fmt.Println("Please input your password:")
	} else {
		warning := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Fprintf(color.Output, "SMEI requires a password to store sensitive information (AudioKinetic and GitHub credentials). %s. Input your password:\n",
			warning("Please note that there is no way to retrieve this password"))
	}

	return passwordLoop()
}

func passwordLoop() error {
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Wrap(err, "could not read a password")
	}

	err = config.SetPassword(secret.String(password))
	if err == config.InvalidPassword {
		color.New(color.FgRed, color.Bold).Println("Invalid password. Please try again.")
		return passwordLoop()
	}
	if err == config.PasswordTooShort {
		color.New(color.FgRed, color.Bold).Println("Password too short. Please try again.")
		return passwordLoop()
	}
	if err != nil {
		return errors.Wrap(err, "could not set the password")
	}

	return nil
}

func askForWwiseAuth() error {
	fmt.Print("SMEI needs credentials to your Audiokinetic/Wwise account. " +
		"If you do not already have one, please navigate to https://www.audiokinetic.com/ and register.\n" +
		"Please input your account email\n")
	//var email string
	//_, err := fmt.Scanln(&email)
	email, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Wrap(err, "could not read the input")
	}
	fmt.Println("Please input your account password: ")
	return wwisePasswordLoop(string(email))
}

func wwisePasswordLoop(email string) error {
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Wrap(err, "could not read a password")
	}

	wwiseClient := client.NewWwiseClient()
	err = wwiseClient.Authenticate(email, string(password))
	if err != nil {
		fmt.Println("Authentication failed. Please try again.")
		return wwisePasswordLoop(email)
	}

	if err != nil {
		return errors.Wrap(err, "could not set the password")
	}

	err = config.SetSecretString(config.WwiseEmail_key, secret.String(email))
	if err != nil {
		return errors.Wrap(err, "could not persist the config change")
	}

	err = config.SetSecretString(config.WwisePassword_key, secret.String(password))
	if err != nil {
		return errors.Wrap(err, "could not persist the config change")
	}

	err = viper.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "could not persist the config change")
	}

	return nil
}
