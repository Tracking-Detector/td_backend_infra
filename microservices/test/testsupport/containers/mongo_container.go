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
	MongoDBContainer struct {
		testcontainers.Container
		//add Config
		Config MongoDBContainerConfig
	}
	//also add options pattern method
	MongoDBContainerOption func(c *MongoDBContainerConfig)

	MongoDBContainerConfig struct {
		ImageTag   string
		MappedPort string
		Database   string
		Host       string
	}
)

func (c *MongoDBContainer) getDns() string {
	return fmt.Sprintf("mongodb://localhost:%s/%s", c.Config.MappedPort, c.Config.Database)
}

func (c *MongoDBContainer) Setenvs() {
	os.Setenv("MONGO_URI", c.getDns())
	os.Setenv("DB_NAME", c.Config.Database)
	os.Setenv("USER_COLLECTION", "users")
	os.Setenv("REQUEST_COLLECTION", "requests")
	os.Setenv("TRAINING_RUNS_COLLECTION", "training-runs")
	os.Setenv("MODELS_COLLECTION", "models")
	os.Setenv("EXPORTER_COLLECTION", "exporter")
	os.Setenv("EXPORTER_RUNS_COLLECTION", "exporter-runs")
}

func NewMongoContainer(ctx context.Context, opts ...MongoDBContainerOption) (*MongoDBContainer, error) {
	const (
		psqlImage = "mongo"
		psqlPort  = "27017"
	)

	config := MongoDBContainerConfig{
		ImageTag: "4.4.18",
		Database: "tracking-detector",
	}
	//handle possible options
	for _, opt := range opts {
		opt(&config)
	}

	containerPort := psqlPort

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Env: map[string]string{
				"MONGO_INITDB_DATABASE": config.Database,
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

	return &MongoDBContainer{
		Container: container,
		Config:    config,
	}, nil
}
