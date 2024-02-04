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
	userRepo := repository.NewMongoUserRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	encryptionService := service.NewEncryptionService()
	userService := service.NewUserService(userRepo, encryptionService)
	userService.InitAdmin(ctx)
	userController := controller.NewUserController(userService)

	userController.Start()
}
