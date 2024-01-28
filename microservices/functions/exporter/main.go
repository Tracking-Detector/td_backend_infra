package exporter

import (
	"context"
	"fmt"

	"tds/shared/configs"
	"tds/shared/consumer"
	"tds/shared/job"
	"tds/shared/queue"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/storage"
)

func Main() {
	fmt.Println("Starting exporter...")
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
	internalExportJob := job.NewInternalExportJob(requestRepo, storageService)
	externalExportJob := job.NewExternalExportJob(requestRepo, storageService)
	exportConsumer := consumer.NewExportMessageConsumer(internalExportJob, externalExportJob, rabbitMqAdapter, exporterService)
	exportConsumer.Consume()
	select {}
}
