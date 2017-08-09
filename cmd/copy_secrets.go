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

	"github.com/dtmistry/swarm-tool/action"
	"github.com/dtmistry/swarm-tool/types"
	"github.com/spf13/cobra"
)

var (
	src, srcCertPath, dest, destCertPath, prefix string
	filters, labels                              = []string{}, []string{}
)

// copySecretsCmd represents the secretsMigrate command
var copySecretsCmd = &cobra.Command{
	Use:   "copy-secrets",
	Short: "Copy secrets between Swarm clusters",
	Long: `Copy secrets from one Swarm cluster to other. 
	For security reasons, this command will only work with a TLS enabled daemon`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(src) == 0 {
			return errors.New("source is required")
		} else if len(dest) == 0 {
			return errors.New("destination is required")
		}
		//TODO add regex check for host
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		source := &types.SwarmConnection{
			Host:     src,
			CertPath: srcCertPath,
		}
		dest := &types.SwarmConnection{
			Host:     dest,
			CertPath: destCertPath,
		}
		err := action.CopySecrets(source, dest, filters, labels, prefix)
		if err != nil {
			return err
		}
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(copySecretsCmd)

	copySecretsCmd.Flags().StringVarP(&src, "source", "s", "", "Source Docker host")
	copySecretsCmd.Flags().StringVarP(&dest, "destination", "d", "", "Destination Docker host")
	copySecretsCmd.Flags().StringVarP(&srcCertPath, "source-cert-path", "", "", "Source Docker TLS cert path")
	copySecretsCmd.Flags().StringVarP(&destCertPath, "destination-cert-path", "", "", "Destination Docker TLS cert path")
	copySecretsCmd.Flags().StringArrayVarP(&filters, "filter", "", nil, "Filters used to copy secrets from source cluster")
	copySecretsCmd.Flags().StringArrayVarP(&labels, "label", "", nil, "Labels added to secret in the target cluster ")
	copySecretsCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Prefix to be added while creating secrets in the target cluster")
	copySecretsCmd.MarkFlagRequired("source")
	copySecretsCmd.MarkFlagRequired("destination")
}
