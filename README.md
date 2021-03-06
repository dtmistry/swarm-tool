[![Build Status](https://travis-ci.org/dtmistry/swarm-tool.svg?branch=master)](https://travis-ci.org/travis-ci/travis-web)
# swarm-tool

A collection of admin tasks for a [Docker Swarm](https://docs.docker.com/engine/swarm/) cluster

## Install

[Download](https://github.com/dtmistry/swarm-tool/releases) the binary (optionally rename) and put it in your path. Execute `help` to verify installation

```bash
$ swarm-tool --help
```

## Commands

By default the tool connects to the local daemon socket. [Environment variables](https://docs.docker.com/engine/reference/commandline/cli/#environment-variables) supported by docker cli can also be used.

To connect to a secure remote daemon socket, use the below environment variables

```bash
$ export DOCKER_HOST=tcp://remote-host:remote-port DOCKER_TLS_VERIFY=1 DOCKER_CERT_PATH=/path/to/certs
```
### rotate-secrets

Updates an existing docker swarm secret

*Usage*

```bash
$ swarm-tool rotate-secrets --secret secret --secret-file=/path/to/updated-secret-data
```
`rotate-secrets` will do the following -

* Check if the `secret` exists
* If there are services which are using this secret...
    * Creates a new `temp_secret` with data from `secret-file`
    * Updates services by removing `secret` and adding `temp_secret`
    * Wait for service updates to converge
    * Updates the `secret` with data from `secret-file`
    * Updates services again. This time removing the `temp_secret` and adding the updated `secret`
    * Wait for service updates to converge
    * Removes the `temp_secret`
* If there are no services which are using this secret...
    * Removes the `secret`
    * Create `secret` with data from the `secret-file`
