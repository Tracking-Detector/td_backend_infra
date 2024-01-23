package service

import (
	"tds/shared/messages"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type IPublishService interface {
	EnqueueTrainingJob(modelId string, exporterId string, reducer string)
	EnqueueExportJob(exporterId string, reducer string, dataset string)
}

type PublishService struct {
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
}

func NewPublishService() *PublishService {
	rabbitConn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // TODO put into .env
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to open Channel: %v", err)
	}
	_, err = rabbitCh.QueueDeclare("exports", true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to declare an exports queue: %v", err)
	}

	_, err = rabbitCh.QueueDeclare("training", true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to declare a training queue: %v", err)
	}

	return &PublishService{
		rabbitConn: rabbitConn,
		rabbitCh:   rabbitCh,
	}
}

func (s *PublishService) EnqueueTrainingJob(modelId string, exporterId string, reducer string) {
	job := messages.NewJob("train_model", []string{modelId, exporterId, reducer})
	message, err := job.Serialize()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Error("Error serializing job.")
		return
	}
	err = s.rabbitCh.Publish("", "training", false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Printf("Failed to publish a message to training queue: %v", err)
	}
}

func (s *PublishService) EnqueueExportJob(exporterId string, reducer string, dataset string) {
	job := messages.NewJob("export", []string{exporterId, reducer, dataset})
	message, err := job.Serialize()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Error("Error serializing job.")
		return
	}
	err = s.rabbitCh.Publish("", "exports", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Printf("Failed to publish a message to training queue: %v", err)
	}
}
