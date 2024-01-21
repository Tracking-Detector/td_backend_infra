package main

import (
	"context"
	"net/http"
	"tds/shared/configs"
	"tds/shared/extractor"
	"tds/shared/models"
	"tds/shared/responses"
	"tds/shared/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var modelCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvModelCollection())
var rabbitConn *amqp.Connection
var rabbitCh *amqp.Channel

func SetupAMQP() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare the queues at startup
	_, err = rabbitCh.QueueDeclare("exports", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare an exports queue: %v", err)
	}

	_, err = rabbitCh.QueueDeclare("training", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a training queue: %v", err)
	}
}

func EnqueueExportJob(exportName string) {
	job := models.NewJob("export", []string{exportName})
	message, err := job.Serialize()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "dispatch",
			"error":   err.Error(),
		}).Error("Error serializing job.")
		return
	}
	err = rabbitCh.Publish("", "exports", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		log.Printf("Failed to publish a message to exports queue: %v", err)
	}
}

func EnqueueTrainingJob(modelName string, dataSetName string) {
	job := models.NewJob("train_model", []string{modelName, dataSetName})
	message, err := job.Serialize()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "dispatch",
			"error":   err.Error(),
		}).Error("Error serializing job.")
		return
	}
	err = rabbitCh.Publish("", "training", false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		log.Printf("Failed to publish a message to training queue: %v", err)
	}
}

func RegisterExportJob(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	extractorName := c.Params("extractorName")
	var ext extractor.Extractor
	for _, ex := range extractor.EXTRACTORS {
		if ex.GetName() == extractorName {
			ext = ex
		}
	}
	if ext.GetName() == "" {
		return c.Status(http.StatusNotFound).JSON(responses.ExportJobStartResponse{
			Status:  http.StatusNotFound,
			Message: "Could not find the extractor you want to trigger.",
		})
	}
	EnqueueExportJob(ext.GetName())
	return c.Status(http.StatusOK).JSON(responses.ExportJobStartResponse{
		Status:  http.StatusOK,
		Message: "The export is queued.",
	})
}

func RegisterTrainingJob(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	modelName := c.Params("modelName")
	dataSetName := c.Params("dataSetName")
	var ext extractor.Extractor
	for _, ex := range extractor.EXTRACTORS {
		if ex.GetName() == dataSetName {
			ext = ex
		}
	}
	if ext.GetName() == "" {
		return c.Status(http.StatusNotFound).JSON(responses.ExportJobStartResponse{
			Status:  http.StatusNotFound,
			Message: "Could not find the extractor you want to trigger.",
		})
	}
	EnqueueTrainingJob(modelName, dataSetName)
	return c.Status(http.StatusOK).JSON(responses.ExportJobStartResponse{
		Status:  http.StatusOK,
		Message: "The training job is queued.",
	})
}

func GetAllPossibleExports(c *fiber.Ctx) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	possibleExports := make([]*fiber.Map, 0)
	for _, ext := range extractor.EXTRACTORS {
		possibleExports = append(possibleExports, &fiber.Map{
			"name":        ext.GetName(),
			"location":    ext.GetFileName(),
			"description": ext.GetDescription(),
		})
	}
	return c.Status(http.StatusOK).JSON(responses.ExportTypesResponse{
		Status: http.StatusOK,
		Data:   possibleExports,
	})
}

func GetAllPossibleModels(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var modelDataValues []models.ModelData
	cursor, _ := modelCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var modelData models.ModelData
		cursor.Decode(&modelData)
		modelDataValues = append(modelDataValues, modelData)
	}

	return c.Status(http.StatusOK).JSON(responses.ModelDataResponse{
		Status: http.StatusOK,
		Data:   modelDataValues,
	})
}

func main() {
	time.Sleep(30 * time.Second)
	SetupAMQP()
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/dispatch/health", utils.GetHealth)
	app.Get("/dispatch/export", GetAllPossibleExports)
	app.Get("/dispatch/model", GetAllPossibleModels)
	app.Post("/dispatch/export/:extractorName", RegisterExportJob)
	app.Post("/dispatch/train/:modelName/run/:dataSetName", RegisterTrainingJob)

	app.Listen(":8081")
}
