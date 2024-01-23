package main

import (
	"context"
	"time"

	"tds/shared/configs"
	"tds/shared/consumer"
	"tds/shared/repository"
	"tds/shared/service"
)

func main() {
	time.Sleep(30 * time.Second)
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	requestRepo := repository.NewMongoRequestRepository(db)
	storageService := service.NewMinIOStorageServic()
	exporterRepo := repository.NewMongoExporterRepository(db)
	exporterService := service.NewExporterService(exporterRepo)
	exportConsumer := consumer.NewExportMessageConsumer(requestRepo, storageService, exporterService)
	exportConsumer.Consume()
	select {}
}
