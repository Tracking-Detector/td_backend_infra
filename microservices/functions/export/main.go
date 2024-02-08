package main

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
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
