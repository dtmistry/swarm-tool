package action

import (
	"bytes"
	"fmt"

	"github.com/dtmistry/swarm-tool/swarm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_PREFIX = "temp"
)

func getSecretName(name, prefix string) string {
	var buf bytes.Buffer
	if len(prefix) == 0 {
		buf.WriteString(DEFAULT_PREFIX)
	} else {
		buf.WriteString(prefix)
	}
	buf.WriteString("_")
	buf.WriteString(name)
	return buf.String()
}

func RotateSecret(secretName, secretFile, prefix string) error {
	c, err := swarm.NewSwarmConnection()
	if err != nil {
		return err
	}
	existingSecret, err := c.FindSecret(secretName)
	if err != nil {
		return err
	}
	if len(existingSecret.ID) == 0 {
		return fmt.Errorf("Secret %s does not exist in this swarm", secretName)
	}

	//Find services with secret
	services, err := c.FindServices(secretName)
	if err != nil {
		return err
	}
	var updateServices bool = false
	if len(services) != 0 {
		updateServices = true
		log.Info("The following services will be updated")
		for _, service := range services {
			log.Info(service.Spec.Name)
		}
	} else {
		log.WithFields(log.Fields{
			"secret": secretName,
		}).Warn("Secret is not attached to any services")
	}
	if updateServices {
		tempSecretName := getSecretName(secretName, prefix)
		tempId, err := c.CreateSecret(tempSecretName, secretFile, nil)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"secret": tempSecretName,
			"ID":     tempId,
		}).Info("Created temp secret")

		tempSecretRef := c.SecretReference(tempSecretName, tempId, secretName)

		for _, service := range services {
			existingSpec := service.Spec
			log.WithFields(log.Fields{
				"service": existingSpec.Name,
			}).Info("Updating service...")
			currentSecrets := existingSpec.TaskTemplate.ContainerSpec.Secrets
			i := 0
			for ; i < len(currentSecrets); i++ {
				if currentSecrets[i].SecretID == existingSecret.ID {
					//Remove this secret
					tempSecretRef.File = currentSecrets[i].File
					break
				}
			}
			currentSecrets = append(currentSecrets[:i], currentSecrets[i+1:]...)
			currentSecrets = append(currentSecrets, tempSecretRef)
			resp, err := c.UpdateService(service.ID, existingSpec, service.Version)
			if err != nil {
				return errors.Wrapf(err, "Unable to update service %s", existingSpec.Name)
			}
			log.WithFields(log.Fields{
				"service":  existingSpec.Name,
				"warnings": resp.Warnings,
			}).Info("Service updated")
		}
	}
	id, err := c.UpdateSecret(existingSecret.ID, secretName, secretFile)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"secret": secretName,
		"ID":     id,
	}).Info("Secret Updated")
	return nil
}
