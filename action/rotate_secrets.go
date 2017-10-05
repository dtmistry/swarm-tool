package action

import (
	"bytes"
	"fmt"

	"github.com/dtmistry/swarm-tool/swarm"
	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_PREFIX = "temp"
)

func GetSecretName(name, prefix string) string {
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
	//Get a connection to swarm
	c, err := swarm.NewSwarmConnection()
	if err != nil {
		return err
	}
	//Check if the secret exists
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
	//Update services
	if len(services) != 0 {
		log.WithFields(log.Fields{
			"secret": secretName,
		}).Info("Services attached to this secret will be updated")

		tempSecretName := GetSecretName(secretName, prefix)
		log.WithFields(log.Fields{
			"temp-secret": tempSecretName,
		}).Info("Creating temp secret")
		tempId, err := c.CreateSecret(tempSecretName, secretFile, nil)
		if err != nil {
			return err
		}
		//Create a temp secret reference with the original target
		tempSecretRef := c.SecretReference(tempSecretName, tempId, secretName)
		//Update the services with temp reference and removing the original secret
		log.WithFields(log.Fields{
			"add-secret":    tempSecretName,
			"remove-secret": secretName,
		}).Info("Updating services with temp secret")
		err = c.UpdateServicesWithSecret(services, tempSecretRef, existingSecret.ID)
		if err != nil {
			return err
		}
		//Update the original secret with new file
		log.WithFields(log.Fields{
			"secret": secretName,
			"file":   secretFile,
		}).Info("Updating original secret with new file")
		updatedSecretId, err := c.UpdateSecret(existingSecret.ID, secretName, secretFile)
		if err != nil {
			return err
		}
		//Create an updated reference
		updatedSecretRef := c.SecretReference(secretName, updatedSecretId, secretName)
		//Find the services again with tempSecretName. Need to do this to get the correct spec.Version object
		updatedServices, err := c.FindServices(tempSecretName)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"add-secret":    secretName,
			"remove-secret": tempSecretName,
		}).Info("Updating services with secret")
		//Update the services with new reference and removing the temp secret
		err = c.UpdateServicesWithSecret(updatedServices, updatedSecretRef, tempId)
		if err != nil {
			return err
		}
		//Delete the temp secret
		log.WithFields(log.Fields{
			"secret": tempSecretName,
		}).Info("Removing temp secret")
		err = c.RemoveSecret(tempId)
		if err != nil {
			return err
		}
	} else {
		log.WithFields(log.Fields{
			"secret": secretName,
		}).Warn("Secret is not attached to any services")
		id, err := c.UpdateSecret(existingSecret.ID, secretName, secretFile)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"secret": secretName,
			"ID":     id,
		}).Info("Secret Updated")
	}
	return nil
}
