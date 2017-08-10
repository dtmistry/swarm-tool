// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
