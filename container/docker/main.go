package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
)

func main() {
	os.Setenv("DOCKER_API_VERSION", "1.44")
	os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s %s\n", ctr.ID, ctr.Image, ctr.Command)

	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range networks {
		fmt.Printf("%s %s %s\n", ctr.ID, ctr.Name, ctr.Driver)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range images {
		fmt.Printf("%s %+v %d\n", ctr.ID, ctr.RepoTags, ctr.Size)
		// for k, v := range ctr.Labels {
		// 	fmt.Printf("%s %+v %s: %s %d\n", ctr.ID, ctr.RepoTags, k, v, ctr.Size)

		// }
	}
}
