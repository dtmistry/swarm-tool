package action

import (
	"context"
	"fmt"

	dockertypes "github.com/docker/docker/api/types"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/dtmistry/swarm-tool/types"
	"github.com/dtmistry/swarm-tool/util"
	"github.com/pkg/errors"
)

func CopySecrets(source, target *types.SwarmConnection) error {
	fmt.Printf("Copying secrets from [%s] to [%s]\n", source.Host, target.Host)

	sourceClient, err := util.NewDockerClient(source.Host, source.CertPath)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to source docker host")
	}

	destClient, err := util.NewDockerClient(target.Host, target.CertPath)
	if err != nil {
		return errors.Wrap(err, "Unable to connect to destination docker host")
	}

	secrets, err := sourceClient.SecretList(context.Background(), dockertypes.SecretListOptions{})
	if err != nil {
		return errors.Wrap(err, "Unable to read secrets from source cluster")
	}

	var secretsToCopy []swarm.SecretSpec

	failedToRead := make(map[string]error)

	for _, secret := range secrets {
		fmt.Printf("Reading secret [%s] from source cluster\n", secret.Spec.Name)
		_, data, err := sourceClient.SecretInspectWithRaw(context.Background(), secret.ID)
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

	if len(failedToRead) != 0 {
		fmt.Println("Unable to inspect the following secrets from the source cluster")
		for secret, err := range failedToRead {
			fmt.Printf("%s : %s", secret, err)
		}
	}

	if len(secretsToCopy) == 0 {
		return errors.New("Unable to read any secrets from the source cluster")
	}

	fmt.Printf("Creating secrets in target cluster\n\n")
	failedToCreate := make(map[string]error)
	for _, secret := range secretsToCopy {
		res, err := destClient.SecretCreate(context.Background(), secret)
		if err != nil {
			failedToCreate[secret.Name] = err
			continue
		}
		fmt.Printf("Created secret [%s] with ID [%s]\n", secret.Name, res.ID)
	}

	if len(failedToCreate) != 0 {
		fmt.Println("Unable to create the following secrets in the target cluster")
		for secret, err := range failedToCreate {
			fmt.Printf("%s : %s", secret, err)
		}
	}

	return nil
}
