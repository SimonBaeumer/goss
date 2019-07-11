package docker_client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	orig_client "github.com/docker/docker/client"
	"github.com/fsouza/go-dockerclient"
)

// DockerClient holds the client which is connected to the docker daemon
type DockerClient struct {
	client *orig_client.Client
}

// Container represents the container structure
type Container struct {
	ID      string
	Running bool
	Names   []string
	Image   string
}

// NewDockerClientFromEnv creates a new client from your current environment
func NewDockerClientFromEnv() *DockerClient {
	client, err := orig_client.NewEnvClient()
	if err != nil {
		panic(err.Error())
	}

	return &DockerClient{
		client: client,
	}
}

// List returns containers by given name
func (d *DockerClient) List(filters map[string][]string) ([]Container, error) {
	opts := docker.ListContainersOptions{Filters: filters}

	containers, err := d.client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: parse.Flag(),
	})
	if err != nil {
		return []Container{}, fmt.Errorf("Error occured: %s", err.Error())
	}

	var result []Container
	for _, container := range containers {
		result = append(result, containerFromAPIContainer(container))
	}
	return result, nil
}

func containerFromAPIContainer(containers docker.APIContainers) Container {
	return Container{
		ID:    containers.ID,
		Image: containers.Image,
		Names: containers.Names,
	}
}
