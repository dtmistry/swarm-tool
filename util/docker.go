package util

import (
	"net/http"
	"path/filepath"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

const (
	DOCKER_API_VERSION = "1.30"
)

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

func NewDockerClient(host, certPath string) (*client.Client, error) {

	var httpClient *http.Client

	if len(certPath) != 0 {
		options := tlsconfig.Options{
			CAFile:             filepath.Join(certPath, "ca.pem"),
			CertFile:           filepath.Join(certPath, "cert.pem"),
			KeyFile:            filepath.Join(certPath, "key.pem"),
			InsecureSkipVerify: false,
		}

		tlsc, err := tlsconfig.Client(options)

		if err != nil {
			return nil, err
		}

		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	}

	return client.NewClient(host, DOCKER_API_VERSION, httpClient, nil)
}
