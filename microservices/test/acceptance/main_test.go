package acceptance

import (
	"context"
	"tds/test/testsupport/containers"
)

type AcceptanceTest struct {
	ctx      context.Context
	cancel   context.CancelFunc
	mongo    *containers.MongoDBContainer
	minio    *containers.MinIOContainer
	rabbitmq *containers.RabbitMQContainer
}

func (t *AcceptanceTest) setupIntegration() {
	t.ctx, t.cancel = context.WithCancel(context.Background())
	t.mongo, _ = containers.NewMongoContainer(t.ctx)
	t.minio, _ = containers.NewMinIOContainer(t.ctx)
	t.rabbitmq, _ = containers.NewRabbitMQContainer(t.ctx)
	t.mongo.Setenvs()
	t.minio.Setenvs()
	t.rabbitmq.Setenvs()
}

func (t *AcceptanceTest) teardownIntegration() {
	t.mongo.Terminate(t.ctx)
	t.minio.Terminate(t.ctx)
	t.rabbitmq.Terminate(t.ctx)
	t.cancel()
}

// func TestMain(m *testing.M) {
// 	t := &AcceptanceTest{}
// 	t.setupIntegration()
// 	code := m.Run()
// 	t.teardownIntegration()
// 	os.Exit(code)
// }
