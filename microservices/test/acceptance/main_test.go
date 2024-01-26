package acceptance

import (
	"context"
	"os"
	"tds/test/testsupport/containers"
	"testing"
	"time"
)

func setupIntegration() {
	// Start docker containers for
	ctx, ctxCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer ctxCancel()
	mongo, _ := containers.NewMongoContainer(ctx)
	minio, _ := containers.NewMinIOContainer(ctx)
	rabbitmq, _ := containers.NewRabbitMQContainer(ctx)
	mongo.Setenvs()
	minio.Setenvs()
	rabbitmq.Setenvs()
}

func TestMain(m *testing.M) {
	setupIntegration()
	code := m.Run()
	os.Exit(code)
}
