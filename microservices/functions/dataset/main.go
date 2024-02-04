package dataset

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/job"
	"tds/shared/repository"
	"tds/shared/service"

	"github.com/robfig/cron/v3"
)

func StartCron(datasetCalculationJob *job.DatasetMetricJob) {
	c := cron.New()
	c.AddFunc("@hourly", func() {
		datasetCalculationJob.Execute()
	})
	c.Start()
}

func Main() {
	ctx := context.TODO()
	requestRepo := repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	requestService := service.NewRequestService(requestRepo)
	datasetRepo := repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	datasetService := service.NewDatasetService(datasetRepo)

	datasetCalculationJob := job.NewDatasetMetricJob(datasetService, requestService)
	datasetController := controller.NewDatasetController(datasetService)

	go StartCron(datasetCalculationJob)
	go datasetController.Start()

	select {}
}
