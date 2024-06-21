package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Muninn",
	Long:  `All software has versions. This is Muninn's`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("v0.7.0")
	},
}
