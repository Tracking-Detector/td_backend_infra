package main

import (
	"context"
	"time"

	"tds/shared/configs"
	"tds/shared/consumer"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/storage"
)

func main() {
	time.Sleep(30 * time.Second)
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	minioClient := configs.ConnectMinio()
	requestRepo := repository.NewMongoRequestRepository(db)
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	storageService := service.NewMinIOStorageService(minioStorageAdapter)
	exporterRepo := repository.NewMongoExporterRepository(db)
	exporterService := service.NewExporterService(exporterRepo)
	exportConsumer := consumer.NewExportMessageConsumer(requestRepo, storageService, exporterService)
	exportConsumer.Consume()
	select {}
}
