package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	ctx := context.TODO()
	trainingrunRepo := repository.NewMongoTrainingRunsRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	modelRepo := repository.NewMongoModelRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	trainingrunService := service.NewTraingingrunService(trainingrunRepo)
	modelService := service.NewModelService(modelRepo, trainingrunService)
	modelController := controller.NewModelController(trainingrunService, modelService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/models/health", utils.GetHealth)
	app.Get("/models", modelController.GetAllModels)
	app.Post("/models", modelController.CreateModel)
	app.Get("/models/:id", modelController.GetModelById)
	app.Delete("/models/:id", modelController.DeleteModelById)
	app.Get("/models/:id/runs", modelController.GetTrainingRunsByModelId)
	app.Delete("/models/:id/runs/:runId", modelController.DeleteRunById)
	app.Get("/models/runs", modelController.GetTrainingRuns)
	app.Listen(":8081")
}
