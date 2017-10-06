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
	Short: "Rotate a secret in a Docker Swarm cluster",
	Long: `Rotate a secret in a Docker Swarm cluster.

rotate-secret will do the following -


* Check if the secret exists
* If there are services which are using this secret...
    * Creates a new temp_secret with data from secret-file
    * Updates services by removing secret and adding temp_secret
    * Wait for service updates to converge
    * Updates the secret with data from secret-file
    * Updates services again. This time removing the temp_secret and adding the updated secret
    * Wait for service updates to converge
    * Removes the temp_secret
* If there are no services which are using this secret...
    * Removes the secret
    * Create secret with data from the secret-file`,
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
	rotateSecretCmd.Flags().StringVarP(&secretFile, "secret-file", "f", "", "File with new secret data")
}
