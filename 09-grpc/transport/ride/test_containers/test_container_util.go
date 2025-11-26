package test_containers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var allRunners sync.Once

type testContainerRunner struct {
	servicePort        int
	name               string
	image              string
	exposedPorts       []string
	env                map[string]string
	cmd                []string
	hostConfigModifier func(*container.HostConfig)
}

func (o testContainerRunner) Run() (testcontainers.Container, error) {
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Name:               o.name,
		Image:              o.image,
		ExposedPorts:       o.exposedPorts,
		Env:                o.env,
		Cmd:                o.cmd,
		WaitingFor:         wait.ForListeningPort(nat.Port(fmt.Sprintf("%d", o.servicePort))).WithStartupTimeout(timeout),
		HostConfigModifier: o.hostConfigModifier,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("[testContainerRunner] %w", err)
	}
	return container, nil
}
