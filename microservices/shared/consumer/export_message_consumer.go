package consumer

import (
	"context"
	"fmt"
	"sync"
	"tds/shared/configs"
	"tds/shared/job"
	"tds/shared/messages"
	"tds/shared/models"
	"tds/shared/queue"
	"tds/shared/service"
	"time"

	log "github.com/sirupsen/logrus"
)

type IConsumer interface {
	Consume()
}

type ExportMessageConsumer struct {
	Wg                sync.WaitGroup
	interExportJob    job.IExportJob
	externalExportJob job.IExportJob
	exportRunService  service.IExportRunService
	queueAdapter      queue.IQueueChannelAdapter
	exporterService   service.IExporterService
	datasetService    service.IDatasetService
}

func NewExportMessageConsumer(interExportJob job.IExportJob, externalExportJob job.IExportJob, exportRunService service.IExportRunService, queueAdapter queue.IQueueChannelAdapter, exporterService service.IExporterService, datasetService service.IDatasetService) *ExportMessageConsumer {
	return &ExportMessageConsumer{
		interExportJob:    interExportJob,
		externalExportJob: externalExportJob,
		exportRunService:  exportRunService,
		exporterService:   exporterService,
		queueAdapter:      queueAdapter,
		datasetService:    datasetService,
		Wg:                sync.WaitGroup{},
	}
}

func (c *ExportMessageConsumer) Consume() {
	fmt.Println("Starting Export ConsumerService...")
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
	datasetId := jobValue.Args[2]

	exporter, err := c.exporterService.FindByID(ctx, exporterId)
	if err != nil || exporter == nil {
		log.Errorf("Exporter does not exist: %v", err)
		return
	}
	dataset, err := c.datasetService.GetDatasetByID(ctx, datasetId)
	if err != nil || dataset == nil {
		log.Errorf("Dataset does not exist: %v", err)
		return
	}
	c.Wg.Add(1)
	go func() {
		start := time.Now()
		run, err := c.exportRunService.Save(ctx, &models.ExportRun{
			ExporterId: exporter.ID,
			Name:       exporter.Name,
			Reducer:    reducer,
			Dataset:    dataset.ID, // TODO what to save in run?
			Start:      start,
		})
		if err != nil {
			log.Fatal("Failed to save export run:", err)
		}
		switch exporter.Type {
		case models.IN_SERVICE:
			run.Metrics = c.interExportJob.Execute(exporter, reducer, dataset.Label)
		case models.JS:
			run.Metrics = c.externalExportJob.Execute(exporter, reducer, dataset.Label)
		}
		run.End = time.Now()
		c.exportRunService.Save(ctx, run)
		if err != nil {
			log.Errorf("Job finished with an error: %v", err)
		}
		defer c.Wg.Done()
	}()
}
