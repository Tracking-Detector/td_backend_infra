package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/repository"
	"tds/shared/service"
)

func main() {
	ctx := context.Background()
	exporterRepo := repository.NewMongoExporterRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	exporterService := service.NewExporterService(exporterRepo)
	// Inits the db with in code extractors
	exporterService.InitInCodeExports(ctx)

	exporterController := controller.NewExtractorController(exporterService)

	exporterController.Start()
}
