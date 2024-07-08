package http

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Container is a testcontainers container
type Container struct {
	testcontainers.Container
	URI string
}

// StartContainer starts a new nginx container
func StartContainer(ctx context.Context) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "nginx:1.23.3",
		ExposedPorts: []string{
			"80",
		},
		WaitingFor: wait.ForHTTP("/").WithPort("80").WithStartupTimeout(10 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	return &Container{
		Container: container,
		URI:       "http://localhost:80",
	}, nil
}
