package config

import (
	"SMEI/lib/ue"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var ConfigDir string
var CacheDir string

var setup = false

type ConfigKey string

const (
	GHClientID_key          = "gh-client-id"
	UEInstallPath_key       = "ue-install-path"
	PreserveUEInstaller_key = "ue-preserve-installer"
	VSInstallPath_key       = "vs-install-path"
	WwiseCacheDir_key       = "cache-dir"
	WwiseEmail_key          = "wwise-email"
	WwisePassword_key       = "wwise-password"
)

func init() {
	setupConfigDir()
	setupCacheDir()
}

func Setup() {
	viper.SetEnvPrefix("SMEI")
	viper.AutomaticEnv()

	setupDefaults()

	viper.AddConfigPath(ConfigDir)
	err := viper.ReadInConfig()
	_, notFound := err.(viper.ConfigFileNotFoundError)
	if err != nil && !notFound {
		log.Fatalf("Could not search for configuration files: %v", err)
	}

	setup = true
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
	viper.SetDefault(UEInstallPath_key, filepath.Join(os.ExpandEnv("$ProgramFiles"), ue.FolderName))
	viper.SetDefault(PreserveUEInstaller_key, false)
	viper.SetDefault(VSInstallPath_key, filepath.Join(os.ExpandEnv("$ProgramFiles"), "Microsoft Visual Studio", "2022", "Community"))
	viper.SetDefault(WwiseCacheDir_key, filepath.Join(CacheDir, "Wwise"))
}
