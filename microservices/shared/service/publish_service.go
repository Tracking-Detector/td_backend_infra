package service

import (
	"fmt"
	"tds/shared/configs"
	"tds/shared/messages"
	"tds/shared/queue"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type IPublishService interface {
	EnqueueTrainingJob(modelId string, exporterId string, reducer string)
	EnqueueExportJob(exporterId string, reducer string, dataset string)
}

type PublishService struct {
	queueAdapter queue.IQueueChannelAdapter
}

func NewPublishService(queueAdapter queue.IQueueChannelAdapter) *PublishService {
	return &PublishService{
		queueAdapter: queueAdapter,
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
	err = s.queueAdapter.Publish("", configs.EnvTrainQueueName(), false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to publish a message to training queue: %v", err)
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
	err = s.queueAdapter.Publish("", configs.EnvExportQueueName(), false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(message),
		DeliveryMode: amqp.Persistent,
	})
	fmt.Println("Published message to export queue")
	if err != nil {
		log.WithFields(log.Fields{
			"service": "PublishService",
			"error":   err.Error(),
		}).Fatalf("Failed to publish a message to training queue: %v", err)
	}
}
