package main

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/job"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/repository"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/robfig/cron/v3"
)

func StartCron(datasetCalculationJob *job.DatasetMetricJob) {
	c := cron.New()
	c.AddFunc("@hourly", func() {
		datasetCalculationJob.Execute()
	})
	c.Start()
}

func main() {
	ctx := context.TODO()
	db := configs.GetDatabase(configs.ConnectDB(ctx))
	requestRepo := repository.NewMongoRequestRepository(db)
	requestService := service.NewRequestService(requestRepo)
	datasetRepo := repository.NewMongoDatasetRepository(db)
	datasetService := service.NewDatasetService(datasetRepo)

	datasetCalculationJob := job.NewDatasetMetricJob(datasetService, requestService)
	datasetController := controller.NewDatasetController(datasetService)

	go StartCron(datasetCalculationJob)
	go datasetController.Start()

	select {}
}
