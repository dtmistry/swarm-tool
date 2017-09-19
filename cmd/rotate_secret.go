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
	"errors"
	"fmt"

	"github.com/dtmistry/swarm-tool/action"
	"github.com/spf13/cobra"
)

var (
	secret, secretFile string
)

// rotateSecretCmd represents the rotateSecret command
var rotateSecretCmd = &cobra.Command{
	Use:   "rotate-secret",
	Short: "Rotate a secret for a docker service",
	Long: `Rotate a secret for a docker services.

rotate-secret will do the following

1) Create a new temp secret with a provided file or string
2) Update a/all sevice/services by adding temp secret and removing the old one
3) Remove the old secret
4) Update the temp secret name with old secret name`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(secret) == 0 {
			return errors.New("secret is required")
		}
		if len(secretFile) == 0 {
			return errors.New("secret-file is required")
		}
		if !PathExists(secretFile) {
			return fmt.Errorf("path %s does not exist", secretFile)
		}
		return action.RotateSecret(secret, secretFile, "")
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(rotateSecretCmd)
	rotateSecretCmd.Flags().StringVarP(&secret, "secret", "s", "", "Secret to rotate")
	rotateSecretCmd.Flags().StringVarP(&secretFile, "secret-file", "f", "", "Secret file to update the secret with")

}
