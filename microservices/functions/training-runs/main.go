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
	trainingrunService := service.NewTraingingrunService(trainingrunRepo)
	trainingrunController := controller.NewTrainingrunController(trainingrunService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/training-runs/health", utils.GetHealth)
	app.Get("/training-runs", trainingrunController.GetTrainingRuns)
	app.Get("/training-runs/:modelName", trainingrunController.GetTrainingRunsByModelName)
	app.Listen(":8081")
}
