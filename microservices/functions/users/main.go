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
	userRepo := repository.NewMongoUserRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	encryptionService := service.NewEncryptionService()
	userService := service.NewUserService(userRepo, encryptionService)
	userService.InitAdmin(ctx)
	userController := controller.NewUserController(userService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/users/health", utils.GetHealth)
	app.Get("/users", userController.GetUsers)
	app.Post("/users", userController.CreateApiUser)
	app.Delete("/users/:Id", userController.DeleteUserByID)
	app.Listen(":8081")
}
