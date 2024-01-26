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
	RabbitMQContainer struct {
		testcontainers.Container
		//add Config
		Config RabbitMQContainerConfig
	}
	//also add options pattern method
	RabbitMQContainerOption func(c *RabbitMQContainerConfig)

	RabbitMQContainerConfig struct {
		ImageTag   string
		MappedPort string
		User       string
		Password   string
		Host       string
	}
)

func (c *RabbitMQContainer) getDns() string {
	return fmt.Sprintf("amqp://%s:%s@localhost:%s/", c.Config.User, c.Config.Password, c.Config.MappedPort)
}

func (c *RabbitMQContainer) Setenvs() {
	os.Setenv("RABBIT_URI", c.getDns())
	os.Setenv("TRAIN_QUEUE", "training")
	os.Setenv("EXPORT_QUEUE", "exports")
}

func NewRabbitMQContainer(ctx context.Context, opts ...RabbitMQContainerOption) (*RabbitMQContainer, error) {
	const (
		psqlImage = "rabbitmq"
		psqlPort  = "5672"
	)

	config := RabbitMQContainerConfig{
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
				"RABBITMQ_DEFAULT_USER": config.User,
				"RABBITMQ_DEFAULT_PASS": config.Password,
			},
			ExposedPorts: []string{
				containerPort,
			},
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

	return &RabbitMQContainer{
		Container: container,
		Config:    config,
	}, nil
}
