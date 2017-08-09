package util

import (
	"net/http"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/fatih/color"
)

const (
	DOCKER_API_VERSION = "1.30"
)

func Err(format string, args ...interface{}) {
	color.Red(format, args...)
}

func Info(format string, args ...interface{}) {
	color.Blue(format, args...)
}

func Warn(format string, args ...interface{}) {
	color.Yellow(format, args...)
}

func NewDockerClient(host, certPath string) (*client.Client, error) {

	httpClient := &http.Client{}

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
