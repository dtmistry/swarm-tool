package util

import (
	"net/http"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

const (
	DOCKER_API_VERSION = "1.303030303030303030303030303030303030303030303030303030303030"
)

func NewDockerClient(host, certPath string) (*client.Client, error) {

	options := tlsconfig.Options{
		CAFile:             filepath.Join(certPath, "ca.pem"),
		CertFile:           filepath.Join(certPath, "cert.pem"),
		KeyFile:            filepath.Join(certPath, "key.pem"),
		InsecureSkipVerify: true,
	}

	tlsc, err := tlsconfig.Client(options)

	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}

	return client.NewClient(host, DOCKER_API_VERSION, httpClient, nil)
}
