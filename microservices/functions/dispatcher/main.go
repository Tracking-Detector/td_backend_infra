package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/queue"
	"tds/shared/repository"
	"tds/shared/service"
	"time"
)

func main() {
	time.Sleep(30 * time.Second)
	ctx := context.Background()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	rabbitCh := configs.ConnectRabbitMQ()

	channelAdapter := queue.NewRabbitMQChannelAdapter(rabbitCh)
	exporterRepo := repository.NewMongoExporterRepository(db)
	modelRepo := repository.NewMongoModelRepository(db)
	trainingRunRepo := repository.NewMongoTrainingRunRepository(db)
	exportRunRepo := repository.NewMongoExportRunRunRepository(db)
	datasetRepo := repository.NewMongoDatasetRepository(db)

	trainingRunService := service.NewTraingingrunService(trainingRunRepo)
	exporterService := service.NewExporterService(exporterRepo)
	exportRunService := service.NewExportRunService(exportRunRepo)
	datasetService := service.NewDatasetService(datasetRepo)
	modelService := service.NewModelService(modelRepo, trainingRunService)
	publisherService := service.NewPublishService(channelAdapter)

	dispatchController := controller.NewDispatchController(exporterService, publisherService, modelService, datasetService, exportRunService)

	dispatchController.Start()
}
