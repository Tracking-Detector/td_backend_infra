package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	time.Sleep(30 * time.Second)
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))

	exporterRepo := repository.NewMongoExporterRepository(db)
	modelRepo := repository.NewMongoModelRepository(db)
	trainingRunRepo := repository.NewMongoTrainingRunsRepository(db)

	trainingRunService := service.NewTraingingrunService(trainingRunRepo)
	exporterService := service.NewExporterService(exporterRepo)
	modelService := service.NewModelService(modelRepo, trainingRunService)
	publisherService := service.NewPublishService()

	dispatchController := controller.NewDispatchController(exporterService, publisherService, modelService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/dispatch/health", utils.GetHealth)
	app.Post("/dispatch/export/:exporterId/:reducer", dispatchController.DispatchExportJob)
	app.Post("/dispatch/train/:modelId/run/:exporterId/:reducer", dispatchController.DispatchTrainingJob)

	app.Listen(":8081")
}
