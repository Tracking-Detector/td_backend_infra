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
	requestRepo := repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	requestService := service.NewRequestService(requestRepo)
	requestController := controller.NewRequestController(requestService)
	requestController.Start()
}
