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
	requestRepo := repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	requestService := service.NewRequestService(requestRepo)
	requestController := controller.NewRequestController(requestService)
	requestController.Start()
}
