package consumer

import (
	"context"
	"sync"
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
	Wg                sync.WaitGroup
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
		Wg:                sync.WaitGroup{},
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

	log.Println("Export ConsumerService started. Waiting for messages...")

	for msg := range msgs {
		c.handleMessage(msg.Body)
	}

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
	c.Wg.Add(1)
	go func() {
		// start := time.Now()
		// metrics := &models.ExportMetrics{}
		switch exporter.Type {
		case models.IN_SERVICE:
			_ = c.interExportJob.Execute(exporter, reducer, dataset)
		case models.JS:
			_ = c.externalExportJob.Execute(exporter, reducer, dataset)
		}
		// end := time.Now()
		if err != nil {
			log.Errorf("Job finished with an error: %v", err)
		}
		defer c.Wg.Done()
	}()
	// TODO write jobs into mongodb

}
