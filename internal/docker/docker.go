package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const appNetwork = "paas"

func RunContainer(appName string, containerName string, image string) (string, error) {
	ctx := context.Background()

	cli, err := client.New(client.FromEnv)

	if err != nil {
		return "", err
	}
	defer cli.Close()

	pullResp, err := cli.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	if err := pullResp.Wait(ctx); err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: containerName,
		Config: &container.Config{
			Image:        image,
			ExposedPorts: network.PortSet{network.MustParsePort("80/tcp"): struct{}{}},
			Labels: map[string]string{
				"traefik.enable": "true",
				fmt.Sprintf("traefik.http.routers.%s.rule", appName):                      fmt.Sprintf("Host(`%s.localhost`)", appName),
				fmt.Sprintf("traefik.http.routers.%s.entrypoints", appName):               "web",
				fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", appName): "80",
			},
		},
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				appNetwork: {},
			},
		},
	})

	if err != nil {
		return "", err
	}

	_, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func DeleteContainer(containerID string) error {
	ctx := context.Background()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	_, _ = cli.ContainerStop(ctx, containerID, client.ContainerStopOptions{})

	_, err = cli.ContainerRemove(
		ctx,
		containerID,
		client.ContainerRemoveOptions{
			Force: true,
		},
	)

	return err
}
