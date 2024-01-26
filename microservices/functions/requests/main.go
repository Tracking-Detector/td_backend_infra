package requests

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

func Main() {
	ctx := context.TODO()
	requestRepo := repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	requestService := service.NewRequestService(requestRepo)
	requestController := controller.NewRequestController(requestService)
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/requests/health", utils.GetHealth)
	app.Get("/requests/:id", requestController.GetRequestById)
	app.Post("/requests", requestController.CreateRequestData)
	app.Post("/requests/multiple", requestController.CreateMultipleRequestData)
	app.Get("/requests", requestController.SearchRequests)
	app.Listen(":8081")
}
