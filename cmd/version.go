package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the swarm-tool version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0.3")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
