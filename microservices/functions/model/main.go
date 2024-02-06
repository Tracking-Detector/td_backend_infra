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
	trainingrunRepo := repository.NewMongoTrainingRunRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	modelRepo := repository.NewMongoModelRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	trainingrunService := service.NewTraingingrunService(trainingrunRepo)
	modelService := service.NewModelService(modelRepo, trainingrunService)
	modelController := controller.NewModelController(trainingrunService, modelService)

	modelController.Start()
}
