package action

import (
	"context"

	dockertypes "github.com/docker/docker/api/types"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/dtmistry/swarm-tool/types"
	"github.com/dtmistry/swarm-tool/util"
	"github.com/pkg/errors"
)

func CopySecrets(source, target *types.SwarmConnection) error {
	util.Info("\nCopying secrets from [%s] to [%s]\n\n", source.Host, target.Host)

	//Create and check source docker client
	srcClient, err := util.NewDockerClient(source.Host, source.CertPath)
	if err != nil {
		return errors.Wrap(err, "Unable to create a client for source docker host")
	}
	_, err = srcClient.Ping(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to connect to source docker host")
	}

	//Create and check tager docker client
	destClient, err := util.NewDockerClient(target.Host, target.CertPath)
	if err != nil {
		return errors.Wrap(err, "Unable to create a client for destination docker host")
	}
	_, err = destClient.Ping(context.Background())
	if err != nil {
		return errors.Wrap(err, "Unable to connect to target docker host")
	}

	//Get a list of all secrets from the source cluster
	secrets, err := srcClient.SecretList(context.Background(), dockertypes.SecretListOptions{})
	if err != nil {
		return errors.Wrap(err, "Unable to read secrets from source cluster")
	}

	var secretsToCopy []swarm.SecretSpec

	failedToRead := make(map[string]error)

	//For all secrets in the list, read the raw data and create a new list with
	// `_copy` as a suffix for the new name
	for _, secret := range secrets {
		util.Info("Reading secret [%s] from source cluster\n", secret.Spec.Name)
		//Get raw secret data for a given secret
		_, data, err := srcClient.SecretInspectWithRaw(context.Background(), secret.ID)
		if err != nil {
			failedToRead[secret.Spec.Name] = err
			continue
		}
		newSecret := swarm.SecretSpec{
			Data:        data,
			Annotations: secret.Spec.Annotations,
		}
		newSecret.Name = secret.Spec.Name + "_copy"
		secretsToCopy = append(secretsToCopy, newSecret)
	}

	//Check if there were errors while reading the secrets
	if len(failedToRead) != 0 {
		util.Warn("Unable to inspect the following secrets from the source cluster")
		for secret, err := range failedToRead {
			util.Warn("%s : %s\n", secret, err)
		}
	}

	//Nothing to copy
	if len(secretsToCopy) == 0 {
		return errors.New("Unable to read any secrets from the source cluster")
	}

	util.Info("\nCreating secrets in target cluster\n\n")
	failedToCreate := make(map[string]error)
	for _, secret := range secretsToCopy {
		//Create the secret in the target cluster
		res, err := destClient.SecretCreate(context.Background(), secret)
		if err != nil {
			failedToCreate[secret.Name] = err
			continue
		}
		util.Info("Created secret [%s] with ID [%s]\n", secret.Name, res.ID)
	}

	if len(failedToCreate) != 0 {
		util.Warn("Unable to create the following secrets in the target cluster")
		for secret, err := range failedToCreate {
			util.Warn("%s : %s\n", secret, err)
		}
	}

	if len(failedToCreate) == len(secretsToCopy) {
		return errors.New("No secrets created in the target cluster")
	}
	return nil
}
