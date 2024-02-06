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
	userRepo := repository.NewMongoUserRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	encryptionService := service.NewEncryptionService()
	userService := service.NewUserService(userRepo, encryptionService)
	userService.InitAdmin(ctx)
	userController := controller.NewUserController(userService)

	userController.Start()
}
