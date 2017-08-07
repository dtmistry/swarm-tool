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

	"github.com/carsdotcom/swarm-tool/action"
	"github.com/carsdotcom/swarm-tool/types"
	"github.com/spf13/cobra"
)

var (
	src, srcCertPath, dest, destCertPath string
)

// secretsMigrateCmd represents the secretsMigrate command
var secretsMigrateCmd = &cobra.Command{
	Use:   "secrets-migrate",
	Short: "Migrate secrets from one Swarm cluster to other",
	Long: `Migrate secrets from one Swarm cluster to other. 
	For security reasons, this command will only work with a TLS enabled daemon`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(src) == 0 {
			return errors.New("source is required")
		} else if len(srcCertPath) == 0 {
			return errors.New("source-cert-path is required")
		} else if len(dest) == 0 {
			return errors.New("destination is required")
		} else if len(destCertPath) == 0 {
			return errors.New("destination-cert-path is required")
		}
		//TODO add regex check for host
		//TODO add file check for cert dirs
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("secrets-migrate called")
		source := &types.SwarmConnection{
			Host:     src,
			CertPath: srcCertPath,
		}
		dest := &types.SwarmConnection{
			Host:     dest,
			CertPath: destCertPath,
		}
		action.MigrateSecrets(source, dest)
	},
}

func init() {
	RootCmd.AddCommand(secretsMigrateCmd)

	secretsMigrateCmd.Flags().StringVarP(&src, "source", "s", "", "Source Docker host")
	secretsMigrateCmd.Flags().StringVarP(&dest, "destination", "d", "", "Destination Docker host")
	secretsMigrateCmd.Flags().StringVarP(&srcCertPath, "source-cert-path", "", "", "Source Docker TLS cert path")
	secretsMigrateCmd.Flags().StringVarP(&destCertPath, "destination-cert-path", "", "", "Destination Docker TLS cert path")
	secretsMigrateCmd.MarkFlagRequired("source")
	secretsMigrateCmd.MarkFlagRequired("destination")
	secretsMigrateCmd.MarkFlagRequired("source-cert-path")
	secretsMigrateCmd.MarkFlagRequired("destination-cert-path")
}
