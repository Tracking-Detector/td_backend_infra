package main

import (
	"github.com/Tracking-Detector/td_backend_infra/dashboard/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/static", "static")
	handlers.NewHomeHandler(app).RegisterHandlers()
	handlers.NewDatasetHandler(app).RegisterHandlers()
	handlers.NewExperimentHandler(app).RegisterHandlers()
	handlers.NewUserHandler(app).RegisterHandlers()
	handlers.NewModelHandler(app).RegisterHandlers()
	handlers.NewExportHandler(app).RegisterHandlers()
	app.Listen(":8081")
}
