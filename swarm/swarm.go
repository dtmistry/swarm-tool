package swarm

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type SwarmConnection struct {
	client *client.Client
}

func NewSwarmConnection() (*SwarmConnection, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return &SwarmConnection{}, err
	}
	return &SwarmConnection{
		client: cli,
	}, nil
}

// Verifies if a Swarm service has a secret attached to it
func (c *SwarmConnection) hasSecret(service swarm.Service, secret string) bool {
	containerSpec := service.Spec.TaskTemplate.ContainerSpec
	secrets := containerSpec.Secrets
	for _, s := range secrets {
		if s.SecretName == secret {
			return true
		}
	}
	return false
}

// Creates a new SecretReference
func (c *SwarmConnection) SecretReference(secretName, secretId, targetName string) *swarm.SecretReference {
	return &swarm.SecretReference{
		SecretName: secretName,
		SecretID:   secretId,
		File: &swarm.SecretReferenceFileTarget{
			Name: targetName,
		},
	}
}

// Finds a secret by name
func (c *SwarmConnection) FindSecret(secretName string) (swarm.Secret, error) {
	secret := swarm.Secret{}
	args := filters.NewArgs()
	args.Add("name", secretName)
	options := types.SecretListOptions{
		Filters: args,
	}
	secrets, err := c.client.SecretList(context.Background(), options)
	if len(secrets) == 0 && err == nil {
		return secret, nil
	}
	if len(secrets) > 1 && err == nil {
		return secret, fmt.Errorf("Multiple secrets found for name: %s", secretName)
	}
	return secrets[0], err
}

func (c *SwarmConnection) UpdateSecret(ID, name, file string) (string, error) {
	err := c.RemoveSecret(ID)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to remove old secret %s", name)
	}
	id, err := c.CreateSecret(name, file, nil)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to update secret %s", name)
	}
	return id, nil
}

// Removes a secret by name
func (c *SwarmConnection) RemoveSecret(id string) error {
	return c.client.SecretRemove(context.Background(), id)
}

// Creates a swarm secret
func (c *SwarmConnection) CreateSecret(secretName, secretFile string, labels map[string]string) (string, error) {
	if len(secretName) == 0 {
		return "", errors.New("secretName cannot be empty")
	}
	data, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to read secret file at %s", secretFile)
	}
	annotations := swarm.Annotations{
		Name:   secretName,
		Labels: labels,
	}
	secretSpec := swarm.SecretSpec{
		Annotations: annotations,
		Data:        data,
	}
	resp, err := c.client.SecretCreate(context.Background(), secretSpec)
	if err != nil {
		return "", errors.Wrap(err, "Error creating secret")
	}
	return resp.ID, nil
}

func (c *SwarmConnection) UpdateService(ID string, spec swarm.ServiceSpec, version swarm.Version) (types.ServiceUpdateResponse, error) {
	return c.client.ServiceUpdate(context.Background(), ID, version, spec, types.ServiceUpdateOptions{})
}

// Returns a list of Swarm services filtered by the provided secret name
func (c *SwarmConnection) FindServices(secretName string) ([]swarm.Service, error) {
	if len(secretName) == 0 {
		return nil, errors.New("secretName cannot be empty")
	}

	services, err := c.client.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "Unable to list services")
	}
	var result []swarm.Service
	for _, s := range services {
		if c.hasSecret(s, secretName) {
			result = append(result, s)
		}
	}
	return result, nil
}
