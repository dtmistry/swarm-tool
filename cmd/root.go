package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "swarm-tool",
	Short: "An admin tool for Docker Swarm cluster/s",
	Long:  `swarm-tool lets you perform some admin workflows to a Docker Swarm cluster. Like rotating and migrating secrets`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		Err("error : %s", err)
		os.Exit(1)
	}
}
