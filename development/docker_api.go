package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	//cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(cli.DaemonHost())

	createdContainer, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
	}, nil, nil, "testing_golang1")

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(createdContainer.ID)

	containers, _ := cli.ContainerList(ctx,
		types.ContainerListOptions{
			//All: true,
			//Filters: filters.NewArgs(
			//   filters.KeyValuePair{
			//       Key: "names", Value: "frosty_heyrovsky",
			//   },
			//),
		})

	fmt.Println("Looking up image")
	for _, c := range containers {
		fmt.Println("hey")
		fmt.Println(c.Names)
	}
}
