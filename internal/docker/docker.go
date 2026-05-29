package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RunContainer(appName string, image string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.%s.rule", appName): fmt.Sprintf("Host(`%s.localhost`)", appName),
			},
		},
		nil,
		nil,
		nil,
		appName,
	)

	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
