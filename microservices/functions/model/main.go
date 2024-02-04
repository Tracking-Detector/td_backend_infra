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
	trainingrunRepo := repository.NewMongoTrainingRunRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	modelRepo := repository.NewMongoModelRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	trainingrunService := service.NewTraingingrunService(trainingrunRepo)
	modelService := service.NewModelService(modelRepo, trainingrunService)
	modelController := controller.NewModelController(trainingrunService, modelService)

	modelController.Start()
}
