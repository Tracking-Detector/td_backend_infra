package test_main

import (
	"context"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

var RabbitMQContainer testcontainers.Container
var MongoDBContainer testcontainers.Container

func TestMain(m *testing.M) {
	runIntegrationTests := os.Getenv("RUN_INTEGRATION_TESTS")
	ctx := context.TODO()
	if runIntegrationTests == "true" {
		// Start RabbitMQ container
		rabbitMQReq := testcontainers.ContainerRequest{
			Image:        "rabbitmq:latest",
			ExposedPorts: []string{"5672"},
		}
		RabbitMQContainer, _ = testcontainers.GenericContainer(
			ctx,
			testcontainers.GenericContainerRequest{
				ContainerRequest: rabbitMQReq,
				Started:          true,
			},
		)
		defer RabbitMQContainer.Terminate(ctx)

		rabbitMQHost, _ := RabbitMQContainer.Host(ctx)
		os.Setenv("RABBITMQ_HOST", rabbitMQHost)
		port, _ := RabbitMQContainer.MappedPort(context.TODO(), "5672")
		os.Setenv("RABBITMQ_PORT", string(port))

		// Start MongoDB container
		mongoDBReq := testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017"},
		}
		MongoDBContainer, _ = testcontainers.GenericContainer(
			ctx,
			testcontainers.GenericContainerRequest{
				ContainerRequest: mongoDBReq,
				Started:          true,
			},
		)
		defer MongoDBContainer.Terminate(ctx)

		mongoDBHost, _ := MongoDBContainer.Host(ctx)
		os.Setenv("MONGODB_HOST", mongoDBHost)
		port, _ = RabbitMQContainer.MappedPort(context.TODO(), "27017")
		os.Setenv("MONGODB_PORT", string(port))
	}

	// Run the tests
	code := m.Run()
	os.Exit(code)
}
