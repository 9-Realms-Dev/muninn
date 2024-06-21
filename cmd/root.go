package cmd

import (
	"github.com/9-Realms-Dev/muninn/internal/tui"
	"github.com/9-Realms-Dev/muninn/internal/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "muninn",
		Short: "File based http client",
		Long:  "defined later",
		Run:   ActivateTui,
	}
)

func ActivateTui(cmd *cobra.Command, args []string) {
	util.Logger.Info("starting tui")
	cd, err := os.Getwd()
	if err != nil {
		util.Logger.Error(err)
	}

	tui.StartTui(cd)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// TODO: Add config file and settings
}

// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := os.UserHomeDir()
// 		cobra.CheckErr(err)

// 		// Search config in home directory with name ".cobra" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigType("yaml")
// 		viper.SetConfigName(".cobra")
// 	}

// 	viper.AutomaticEnv()

// 	if err := viper.ReadInConfig(); err == nil {
// 		util.Logger.Info("Using config file:", viper.ConfigFileUsed())
// 	}
// }
