package main

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/consumer"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/job"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/queue"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/storage"
)

func main() {
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	minioClient := configs.ConnectMinio()
	rabbitMqChannel := configs.ConnectRabbitMQ()
	rabbitMqAdapter := queue.NewRabbitMQChannelAdapter(rabbitMqChannel)
	requestRepo := repository.NewMongoRequestRepository(db)
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	storageService := service.NewMinIOStorageService(minioStorageAdapter)
	storageService.VerifyBucketExists(ctx, configs.EnvExtractorBucketName())
	storageService.VerifyBucketExists(ctx, configs.EnvModelBucketName())
	storageService.VerifyBucketExists(ctx, configs.EnvExportBucketName())
	exporterRepo := repository.NewMongoExporterRepository(db)
	exporterService := service.NewExporterService(exporterRepo)
	datasetRepo := repository.NewMongoDatasetRepository(db)
	datasetService := service.NewDatasetService(datasetRepo, requestRepo)
	exportRunRepo := repository.NewMongoExportRunRunRepository(db)
	exportRunService := service.NewExportRunService(exportRunRepo)
	internalExportJob := job.NewInternalExportJob(requestRepo, storageService)
	externalExportJob := job.NewExternalExportJob(requestRepo, storageService)
	exportConsumer := consumer.NewExportMessageConsumer(internalExportJob, externalExportJob, exportRunService, rabbitMqAdapter, exporterService, datasetService)
	exportConsumer.Consume()
	select {}
}
