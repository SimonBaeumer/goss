package system

import (
	"fmt"
	"github.com/SimonBaeumer/goss/internal/docker-client"
	"github.com/SimonBaeumer/goss/util"
)

type Docker interface {
	Image() string
	Running() bool
	Count() int
}

type DefDocker struct {
	ImageName string
}

func createDockerClient() {
	client := docker_client.NewDockerClientFromEnv()
	containers, _ := client.List(map[string][]string{"name": {"confident_wilson"}})

	fmt.Println("Looking up image")
	for _, c := range containers {
		fmt.Println("hey")
		fmt.Println(c.Names)
	}
}

func (d DefDocker) Image() string {
	return d.ImageName
}

func (d DefDocker) Running() bool {
	return true
}

func (d DefDocker) Count() int {
	return 0
}

func NewDefDocker(name string, system *System, config util.Config) Docker {
	return DefDocker{
		ImageName: name,
	}
}
