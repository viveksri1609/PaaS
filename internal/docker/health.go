package docker

import (
	"context"
	"errors"

	"github.com/moby/moby/client"
)

func IsContainerRunning(containerID string) (bool, error) {
	if containerID == "" {
		return false, errors.New("container id is empty")
	}

	ctx := context.Background()

	cli, err := client.New(client.FromEnv)

	if err != nil {
		return false, err
	}
	defer cli.Close()

	containerInfo, err := cli.ContainerInspect(
		ctx,
		containerID,
		client.ContainerInspectOptions{},
	)

	if err != nil {
		return false, err
	}

	if containerInfo.Container.State == nil {
		return false, errors.New("container state is unavailable")
	}

	return containerInfo.Container.State.Running, nil
}
