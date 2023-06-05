package credentials

import (
	"SMEI/config"
	"SMEI/lib/colors"
	"SMEI/lib/secret"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/mircearoata/wwise-cli/lib/wwise/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

type WwiseAuth struct {
	Email    secret.String
	Password secret.String
}

func AskForPassword() error {
	if config.HasLoggedInBefore() {
		colors.RequestColor.Println("If you forgot your password, delete config.yml in '%APPDATA%\\SMEI\\'\nPlease input your password (input is obscured):")
	} else {
		warning := colors.WarningColor.SprintFunc()
		colors.RequestColor.Fprintf(color.Output, "SMEI requires a password to store sensitive information (AudioKinetic and GitHub credentials). %s Create a password (input is obscured):\n",
			warning("Please note that there is no way to retrieve this password."))
	}

	return passwordLoop()
}

func passwordLoop() error {
	password := []byte{}
	err := error(nil)
	if viper.GetBool(config.DeveloperMode_key) {
		colors.SequenceColor.Println("SMEI developer mode enabled. Using default SMEI password for testing.")
		password = []byte("FrenchFeyko")
	} else {
		password, err = terminal.ReadPassword(int(os.Stdin.Fd()))
	}

	if err != nil {
		return errors.Wrap(err, "could not read a password")
	}

	err = config.SetPassword(secret.String(password))
	if err == config.InvalidPassword {
		colors.ErrorColor.Println("Invalid password. Please try again.")
		return passwordLoop()
	}
	if err == config.PasswordTooShort {
		colors.ErrorColor.Println("Password too short. Please try again.")
		return passwordLoop()
	}
	if err != nil {
		return errors.Wrap(err, "could not set the password")
	}

	return nil
}

func askForWwiseAuth() error {
	colors.RequestColor.Print("SMEI needs credentials to your Audiokinetic/Wwise account. " +
		"If you do not already have one, please navigate to https://www.audiokinetic.com/ and register.\n" +
		"Please input your account email (input is obscured):\n")
	email, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Wrap(err, "could not read the input")
	}
	colors.RequestColor.Println("Please input your account password (input is obscured): ")
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

func GetWwiseCredentials() (*WwiseAuth, error) {
	if !viper.IsSet(config.WwiseEmail_key) {
		err := askForWwiseAuth()
		if err != nil {
			log.Panicf("Could not log in with Wwise: %v", err)
		}
	}

	wwiseEmail, err := config.GetSecretString(config.WwiseEmail_key)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get the Wwise email")
	}

	wwisePassword, err := config.GetSecretString(config.WwisePassword_key)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get the Wwise password")
	}

	return &WwiseAuth{
		Email:    wwiseEmail,
		Password: wwisePassword,
	}, nil
}
