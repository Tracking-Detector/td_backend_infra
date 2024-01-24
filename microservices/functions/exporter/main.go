package main

import (
	"context"
	"time"

	"tds/shared/configs"
	"tds/shared/consumer"
	"tds/shared/job"
	"tds/shared/queue"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/storage"
)

func main() {
	time.Sleep(30 * time.Second)
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	minioClient := configs.ConnectMinio()
	rabbitMqChannel := configs.ConnectRabbitMQ()
	rabbitMqAdapter := queue.NewRabbitMQChannelAdapter(rabbitMqChannel)
	requestRepo := repository.NewMongoRequestRepository(db)
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	storageService := service.NewMinIOStorageService(minioStorageAdapter)
	exporterRepo := repository.NewMongoExporterRepository(db)
	exporterService := service.NewExporterService(exporterRepo)
	internalExportJob := job.NewInternalExportJob(requestRepo, storageService)
	externalExportJob := job.NewExternalExportJob(requestRepo, storageService)
	exportConsumer := consumer.NewExportMessageConsumer(internalExportJob, externalExportJob, rabbitMqAdapter, requestRepo, storageService, exporterService)
	exportConsumer.Consume()
	select {}
}
