package containers

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type (
	MinIOContainer struct {
		testcontainers.Container
		//add Config
		Config MinIOContainerConfig
	}
	//also add options pattern method
	MinIOContainerOption func(c *MinIOContainerConfig)

	MinIOContainerConfig struct {
		ImageTag   string
		MappedPort string
		User       string
		Password   string
		Host       string
	}
)

func (c *MinIOContainer) getDns() string {
	return fmt.Sprintf("localhost:%s", c.Config.MappedPort)
}

func (c *MinIOContainer) Setenvs() {
	os.Setenv("MINIO_URI", c.getDns())
	os.Setenv("MINIO_ACCESS_KEY", c.Config.User)
	os.Setenv("MINIO_PRIVATE_KEY", c.Config.Password)
	os.Setenv("EXPORT_BUCKET_NAME", "exports")
	os.Setenv("MODEL_BUCKET_NAME", "models")
	os.Setenv("EXTRACTOR_BUCKET_NAME", "extractors")
}

func NewMinIOContainer(ctx context.Context, opts ...MinIOContainerOption) (*MinIOContainer, error) {
	const (
		psqlImage = "minio/minio"
		psqlPort  = "9000"
	)

	config := MinIOContainerConfig{
		ImageTag: "latest",
		User:     "adminadmin",
		Password: "adminadminadmin",
	}
	//handle possible options
	for _, opt := range opts {
		opt(&config)
	}

	containerPort := psqlPort

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"MINIO_ROOT_USER":     config.User,
				"MINIO_ROOT_PASSWORD": config.Password,
			},
			ExposedPorts: []string{
				containerPort,
			},
			Cmd:        []string{"server", "/data", "--console-address", ":9001"},
			Image:      fmt.Sprintf("%s:%s", psqlImage, config.ImageTag),
			WaitingFor: wait.ForListeningPort(nat.Port(containerPort)),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getting request provider: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host for: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(containerPort))
	if err != nil {
		return nil, fmt.Errorf("getting mapped port for (%s): %w", containerPort, err)
	}
	config.MappedPort = mappedPort.Port()
	config.Host = host

	return &MinIOContainer{
		Container: container,
		Config:    config,
	}, nil
}
