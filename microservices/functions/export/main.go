package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"time"

	"strings"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestDataCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvRequestCollection())

func startExportConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the exports queue
	_, err = ch.QueueDeclare(
		"exports", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exports queue: %v", err)
	}

	msgs, err := ch.Consume(
		"exports", // queue name
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		job, err := models.DeserializeJob(string(msg.Body))
		if err != nil {
			log.Errorf("Failed to deserialize job: %v", err)
			continue
		}

		exportName := job.Args[0]
		var ext extractor.Extractor
		for _, ex := range extractor.EXTRACTORS {
			if ex.GetName() == exportName {
				ext = ex
			}
		}
		if ext.GetName() == "" {
			log.WithFields(log.Fields{
				"service": "export",
			}).Error("Error identifing export job: ", exportName[0])
			return
		}
		RunDataExport(ext)

	}
}

func RunDataExport(fe extractor.Extractor) {
	pr, pw := io.Pipe()

	gzipWriter := gzip.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer gzipWriter.Close()

		cursor, err := requestDataCollection.Find(context.Background(), bson.M{})
		if err != nil {
			log.WithFields(log.Fields{
				"service": "export",
				"error":   err.Error(),
			}).Error("Failed to query MongoDB collection.")
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var doc models.RequestData
			if err := cursor.Decode(&doc); err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Failed to decode MongoDB document.")
				continue
			}
			encoded, encodeErr := fe.Encode(doc)
			if encodeErr != nil {
				continue
			}
			arr, err := json.Marshal(encoded)
			if err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Could not convert int[] to string.")
				continue
			}
			data := strings.Trim(string(arr), "[]") + "\n"

			if _, err := gzipWriter.Write([]byte(data)); err != nil {
				log.WithFields(log.Fields{
					"service": "export",
					"error":   err.Error(),
				}).Error("Failed to write to gzip writer.")
				break
			}
		}

		if cursor.Err() != nil {
			log.WithFields(log.Fields{
				"service": "export",
				"error":   cursor.Err().Error(),
			}).Error("Error occurred while iterating MongoDB cursor.")
		}
		log.WithFields(log.Fields{
			"service": "export",
		}).Info("Data compression and upload for " + fe.GetName() + " completed successfully.")
	}()
	_, putErr := configs.MINIO.PutObject(context.Background(), configs.EnvExportBucketName(), fe.GetFileName(), pr, -1, minio.PutObjectOptions{
		ContentType: "application/gzip",
	})
	if putErr != nil {
		log.WithFields(log.Fields{
			"service": "export",
			"error":   putErr.Error(),
		}).Error("Failed to upload data to MinIO.")
	}

}

func main() {
	time.Sleep(30 * time.Second)
	configs.VerifyBucketExists(context.Background(), configs.MINIO, configs.EnvExportBucketName())
	go startExportConsumer()
	select {}
}
