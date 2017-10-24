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
	log "github.com/sirupsen/logrus"
)

// Defines a connection to a Docker swarm cluster
type SwarmConnection struct {
	client *client.Client
}

// Creates a new SwarmConnection by initializing the docker api client from enviornment defaults
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
		for _, s := range secrets {
			if s.Spec.Name == secretName {
				return s, nil
			}
		}
	}
	return secrets[0], err
}

// Updates a secret by removing it and creating it again with the file provided
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

//Waits for the service to converge
func (c *SwarmConnection) waitForServiceUpdate(ctx context.Context, serviceID string) error {
	for {
		service, _, err := c.client.ServiceInspectWithRaw(ctx, serviceID, types.ServiceInspectOptions{})
		if err != nil {
			return err
		}
		if service.UpdateStatus != nil {
			switch service.UpdateStatus.State {
			case swarm.UpdateStateCompleted:
				return nil
			case swarm.UpdateStatePaused:
				return fmt.Errorf("service update paused: %s", service.UpdateStatus.Message)
			case swarm.UpdateStateRollbackCompleted:
				return fmt.Errorf("service rolled back: %s", service.UpdateStatus.Message)
			case swarm.UpdateStateRollbackPaused:
				return fmt.Errorf("service rollback paused: %s", service.UpdateStatus.Message)
			}
		}
	}
}

//TODO Aggregate errors
//Updates the supplied []swarm.Service with the new secret reference
func (c *SwarmConnection) UpdateServicesWithSecret(services []swarm.Service, newSecretRef *swarm.SecretReference, secretIdToReplace string) error {
	for _, service := range services {
		ctx := context.Background()
		existingSpec := service.Spec
		log.WithFields(log.Fields{
			"service": existingSpec.Name,
		}).Info("Updating service")
		currentSecrets := existingSpec.TaskTemplate.ContainerSpec.Secrets
		i := 0
		for ; i < len(currentSecrets); i++ {
			if currentSecrets[i].SecretID == secretIdToReplace {
				//Remove this secret
				newSecretRef.File = currentSecrets[i].File
				break
			}
		}
		currentSecrets = append(currentSecrets[:i], currentSecrets[i+1:]...)
		currentSecrets = append(currentSecrets, newSecretRef)
		resp, err := c.UpdateService(ctx, service.ID, existingSpec, service.Version)
		if err != nil {
			return errors.Wrapf(err, "Unable to update service %s", existingSpec.Name)
		}
		//Wait for service update to converge
		err = c.waitForServiceUpdate(ctx, service.ID)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"service":  existingSpec.Name,
			"warnings": resp.Warnings,
		}).Info("Service updated")
	}
	return nil
}

// Updates a service
func (c *SwarmConnection) UpdateService(ctx context.Context, ID string, spec swarm.ServiceSpec, version swarm.Version) (types.ServiceUpdateResponse, error) {
	return c.client.ServiceUpdate(ctx, ID, version, spec, types.ServiceUpdateOptions{})
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
