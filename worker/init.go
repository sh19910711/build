package worker

import (
	"github.com/docker/docker/client"
	"os"
)

const DEFAULT_DOCKER_ENDPOINT = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "v1.18"

var dockerClient *client.Client

func init() {
	headers := map[string]string{
		"User-Agent": "cs-build",
	}

	// $ export DOCKER_HOST=tcp://dockerhost:2375
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = DEFAULT_DOCKER_ENDPOINT
	}

	if cli, err := client.NewClient(endpoint, DOCKER_API_VERSION, nil, headers); err != nil {
		panic(err)
	} else {
		dockerClient = cli
	}
}
