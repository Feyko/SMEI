package cmd

import (
	"SMEI/config"
	"SMEI/lib/project"
	"SMEI/lib/ue"
	"SMEI/lib/vs"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

		err := viper.BindPFlags(flags)
		if err != nil {
			log.Fatalf("Could not bind the CLI flags to the configuration system: %v", err)
		}

		viper.Set("config-dir", filepath.Join(config.CacheDir, "Wwise"))

		local := viper.GetBool("local")

		target := getStringFlag(flags, "target")

		fmt.Println("Installing the Unreal Engine")
		err = ue.Install(target, local)
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

func getStringFlag(flags *pflag.FlagSet, name string) string {
	flag, err := flags.GetString(name)
	if err != nil {
		log.Fatalf("Could not get the '%v' flag: %v", name, err)
	}

	if flag == "" {
		log.Fatalf("Config value %v was not set", name)
	}

	return flag
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

	viper.SetEnvPrefix("SMEI")
	viper.AutomaticEnv()

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Could not get the user's configuration directory: %v", err)
	}
	config.ConfigDir = filepath.Join(configDir, "SMEI")

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Could not get the user's cache directory: %v", err)
	}
	config.CacheDir = filepath.Join(cacheDir, "SMEI")

	viper.AddConfigPath(config.ConfigDir)
	err = viper.ReadInConfig()
	_, notFound := err.(viper.ConfigFileNotFoundError)
	if err != nil && !notFound {
		log.Fatalf("Could not search for configuration files: %v", err)
	}
}
