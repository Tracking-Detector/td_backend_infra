package main

import (
	"context"

	"time"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/queue"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
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
