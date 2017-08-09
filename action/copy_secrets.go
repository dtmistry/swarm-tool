package action

import (
	"context"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	swarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/dtmistry/swarm-tool/types"
	"github.com/dtmistry/swarm-tool/util"
	"github.com/pkg/errors"
)

func GetMap(flags []string) (map[string]string, error) {
	args := make(map[string]string)
	if len(flags) == 0 {
		return args, nil
	}
	for i := range flags {
		if !strings.Contains(flags[i], "=") {
			return args, errors.New("bad format of labels (expected name=value)")
		} else {
			parts := strings.SplitN(flags[i], "=", 2)
			name := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			args[name] = value
		}
	}
	return args, nil
}

func GetArgs(flags []string) (filters.Args, error) {
	var (
		args = filters.NewArgs()
		err  error
	)
	for i := range flags {
		args, err = filters.ParseFlag(flags[i], args)
		if err != nil {
			return args, err
		}
	}
	return args, nil
}

func readSecrets(client *client.Client, secrets []swarm.Secret, prefix string, createArgs map[string]string) ([]swarm.SecretSpec, map[string]error) {
	var secretsToCopy []swarm.SecretSpec
	failedToRead := make(map[string]error)

	//For all secrets in the list, read the raw data and create a new list with
	// `_copy` as a suffix for the new name
	for _, secret := range secrets {
		util.Info("Reading secret [%s] from source cluster\n", secret.Spec.Name)
		//Get raw secret data for a given secret
		_, data, err := client.SecretInspectWithRaw(context.Background(), secret.ID)
		if err != nil {
			failedToRead[secret.Spec.Name] = err
			continue
		}
		newSecret := swarm.SecretSpec{
			Data:        data,
			Annotations: secret.Spec.Annotations,
		}
		if len(prefix) == 0 {
			newSecret.Name = secret.Spec.Name
		} else {
			newSecret.Name = prefix + "_" + secret.Spec.Name
		}
		newSecret.Labels = createArgs
		secretsToCopy = append(secretsToCopy, newSecret)
	}
	return secretsToCopy, failedToRead
}

func createSecrets(client *client.Client, secretsToCopy []swarm.SecretSpec) map[string]error {
	util.Info("\nCreating secrets in target cluster\n\n")
	failedToCreate := make(map[string]error)
	for _, secret := range secretsToCopy {
		//Create the secret in the target cluster
		res, err := client.SecretCreate(context.Background(), secret)
		if err != nil {
			failedToCreate[secret.Name] = err
			continue
		}
		util.Info("Created secret [%s] with ID [%s]\n", secret.Name, res.ID)
	}
	return failedToCreate
}

//TODO break the huge method
func CopySecrets(source, target *types.SwarmConnection, filters, labels []string, prefix string) error {
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

	filterArgs, err := GetArgs(filters)
	if err != nil {
		return errors.Wrap(err, "Unable to parse filter labels")
	}

	createArgs, err := GetMap(labels)
	if err != nil {
		return errors.Wrap(err, "Unable to parse labels")
	}

	//Add copied-from label. copied-from=host
	createArgs["copied-from"] = source.Host

	//Get a list of all secrets from the source cluster
	secrets, err := srcClient.SecretList(context.Background(), dockertypes.SecretListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return errors.Wrap(err, "Unable to read secrets from source cluster")
	}

	secretsToCopy, failedToRead := readSecrets(srcClient, secrets, prefix, createArgs)
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

	failedToCreate := createSecrets(destClient, secretsToCopy)

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