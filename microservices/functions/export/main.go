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
	exporterRepo := repository.NewMongoExporterRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	exporterService := service.NewExporterService(exporterRepo)
	// Inits the db with in code extractors
	exporterService.InitInCodeExports(ctx)

	exporterController := controller.NewExtractorController(exporterService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/export/health", utils.GetHealth)
	app.Get("/export", exporterController.GetAllExporter)
	app.Listen(":8081")
}
