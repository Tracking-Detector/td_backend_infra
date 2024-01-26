package configs

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"service": "Setup",
	}).Info("Successfully connected to MongoDB.")
	return client
}

func ConnectMinio() *minio.Client {
	minioClient, err := minio.New(EnvMinIoURI(), &minio.Options{
		Creds:  credentials.NewStaticV4(EnvMinIoAccessKey(), EnvMinIoPrivateKey(), ""),
		Secure: false,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "Setup",
		}).Fatalln(err)
	}
	log.WithFields(log.Fields{
		"service": "Setup",
	}).Info("Successfully connected to MinIO.")
	return minioClient
}

func ConnectRabbitMQ() *amqp.Channel {
	rabbitConn, err := amqp.Dial(EnvMQURI())
	if err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatalf("Failed to open Channel: %v", err)
	}
	_, err = rabbitCh.QueueDeclare(EnvExportQueueName(), true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatalf("Failed to declare an exports queue: %v", err)
	}

	_, err = rabbitCh.QueueDeclare(EnvTrainQueueName(), true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatalf("Failed to declare a training queue: %v", err)
	}
	return rabbitCh
}

func GetDatabase(client *mongo.Client) *mongo.Database {
	db := client.Database(EnvDBName())
	return db
}
