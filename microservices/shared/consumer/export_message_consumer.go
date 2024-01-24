package consumer

import (
	"context"
	"tds/shared/configs"
	"tds/shared/job"
	"tds/shared/messages"
	"tds/shared/models"
	"tds/shared/queue"
	"tds/shared/service"

	log "github.com/sirupsen/logrus"
)

type IConsumer interface {
	Consume()
}

type ExportMessageConsumer struct {
	interExportJob    job.IExportJob
	externalExportJob job.IExportJob
	queueAdapter      queue.IQueueChannelAdapter
	exporterService   service.IExporterService
}

func NewExportMessageConsumer(interExportJob job.IExportJob, externalExportJob job.IExportJob, queueAdapter queue.IQueueChannelAdapter, exporterService service.IExporterService) *ExportMessageConsumer {
	return &ExportMessageConsumer{
		interExportJob:    interExportJob,
		externalExportJob: externalExportJob,
		exporterService:   exporterService,
		queueAdapter:      queueAdapter,
	}
}

func (c *ExportMessageConsumer) Consume() {
	msgs, err := c.queueAdapter.Consume(
		configs.EnvExportQueueName(), // queue name
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
	// TODO how to make this better for testing
	// stopChan := make(chan os.Signal, 1)
	// signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Export ConsumerService started. Waiting for messages...")

	for msg := range msgs {
		c.handleMessage(msg.Body)
	}

	// <-stopChan
	log.Println("Shutting down Export ConsumerService.")
}

func (c *ExportMessageConsumer) handleMessage(msg []byte) {
	ctx := context.TODO()
	jobValue, err := messages.DeserializeJob(string(msg))
	if err != nil {
		log.Errorf("Failed to deserialize job: %v", err)
		return
	}
	exporterId := jobValue.Args[0]
	reducer := jobValue.Args[1]
	dataset := jobValue.Args[2]

	exporter, err := c.exporterService.FindByID(ctx, exporterId)

	if err != nil || exporter == nil {
		log.Errorf("Exporter does not exist: %v", err)
		return
	}
	go func() {
		switch exporter.Type {
		case models.IN_SERVICE:
			err = c.interExportJob.Execute(exporter, reducer, dataset)
		case models.JS:
			err = c.externalExportJob.Execute(exporter, reducer, dataset)
		}
		if err != nil {
			log.Errorf("Job finished with an error: %v", err)
		}
	}()
	// TODO write jobs into mongodb

}
