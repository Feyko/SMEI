package config

import (
	"SMEI/lib/crypt"
	"SMEI/lib/secret"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var ConfigDir string
var CacheDir string

const UEFolderName = "Unreal Engine - CSS"

var password secret.String

const passCheck = "SMEI"
const MinPasswordLength = 8

const (
	GHClientID_key          = "gh-client-id"
	UEInstallPath_key       = "ue-install-path"
	PreserveUEInstaller_key = "ue-preserve-installer"
	VSInstallPath_key       = "vs-install-path"
	WwiseCacheDir_key       = "cache-dir"
	WwiseEmail_key          = "wwise-email"
	WwisePassword_key       = "wwise-password"
	PassCheck_key           = "pass-check"
)

func init() {
	setupConfigDir()
	setupCacheDir()
}

func Setup() {
	viper.SetEnvPrefix("SMEI")
	viper.AutomaticEnv()

	setupDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir)
	err := viper.ReadInConfig()
	_, notFound := err.(viper.ConfigFileNotFoundError)
	if notFound {
		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatalf("Could not write the default config: %v", err)
		}
	}
	if err != nil && !notFound {
		log.Fatalf("Could not search for configuration files: %v", err)
	}
}

func setupConfigDir() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Panicf("Could not get the user's configuration directory: %v", err)
	}
	ConfigDir = filepath.Join(configDir, "SMEI")
}

func setupCacheDir() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Panicf("Could not get the user's cache directory: %v", err)
	}
	CacheDir = filepath.Join(cacheDir, "SMEI")
}

func setupDefaults() {
	viper.SetDefault(GHClientID_key, "0e4260b720ae65240864")
	viper.SetDefault(UEInstallPath_key, filepath.Join(os.ExpandEnv("$ProgramFiles"), UEFolderName))
	viper.SetDefault(PreserveUEInstaller_key, true)
	viper.SetDefault(VSInstallPath_key, filepath.Join(os.ExpandEnv("$ProgramFiles"), "Microsoft Visual Studio", "2022", "Community"))
	viper.SetDefault(WwiseCacheDir_key, filepath.Join(CacheDir, "Wwise"))
}

func SetPassword(newPassword secret.String) error {
	if len(newPassword) < MinPasswordLength {
		return PasswordTooShortError{}
	}

	if !HasLoggedInBefore() {
		encrypted, err := crypt.Encrypt(string(newPassword), passCheck)
		if err != nil {
			return errors.Wrap(err, "could not encrypt the pass check")
		}

		viper.Set(PassCheck_key, encrypted)
		err = viper.WriteConfig()
	}

	decrypted, err := crypt.Decrypt(string(newPassword), viper.GetString(PassCheck_key))
	if err != nil {
		return errors.Wrap(err, "could not decrypt the pass check")
	}

	if decrypted != passCheck {
		return InvalidPasswordError{}
	}

	err = viper.WriteConfig()

	if err != nil {
		return errors.Wrap(err, "could not persist the config changes")
	}

	password = newPassword

	return nil
}

var PasswordTooShort = PasswordTooShortError{}

type PasswordTooShortError struct{}

func (e PasswordTooShortError) Error() string {
	return fmt.Sprintf("password is too short (<%v)", MinPasswordLength)
}

func HasPassword() bool {
	return password != ""
}

func HasLoggedInBefore() bool {
	return viper.IsSet(PassCheck_key)
}

var InvalidPassword = InvalidPasswordError{}

type InvalidPasswordError struct{}

func (e InvalidPasswordError) Error() string {
	return "invalid password"
}

func SetSecretString(key string, str secret.String) error {
	if !HasPassword() {
		return MissingPasswordError{}
	}

	encrypted, err := crypt.Encrypt(string(password), string(str))
	if err != nil {
		return errors.Wrap(err, "could not encrypt")
	}

	viper.Set(key, encrypted)

	return nil
}

var MissingPassword = MissingPasswordError{}

type MissingPasswordError struct{}

func (e MissingPasswordError) Error() string {
	return "missing password"
}

func GetSecretString(key string) (secret.String, error) {
	if !HasPassword() {
		return "", MissingPasswordError{}
	}

	str := viper.GetString(key)
	if str == "" {
		return "", nil
	}

	decrypted, err := crypt.Decrypt(string(password), str)
	if err != nil {
		return "", errors.Wrap(err, "could not decrypt")
	}

	return secret.String(decrypted), nil
}
