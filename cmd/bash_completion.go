package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bashCompletionCmd represents the bashCompletion command
var bashCompletionCmd = &cobra.Command{
	Use:   "bash-completion",
	Short: "Generate bash completion for the swarm-tool project.",
	Long:  `Generate bash completion for the swarm-tool project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := RootCmd.GenBashCompletionFile("contrib/completion/bash/swarm-tool")
		if err != nil {
			return fmt.Errorf("Unable to generate bash completeion %s", err)
		}
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(bashCompletionCmd)
}
