package dataset

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/job"
	"tds/shared/repository"
	"tds/shared/service"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/robfig/cron/v3"
)

func StartServer(datasetController *controller.DatasetController) {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/datasets/health", utils.GetHealth)
	app.Post("/datasets", datasetController.CreateDataset)
	app.Get("/datasets", datasetController.GetAllDatasets)
	app.Listen(":8081")
}

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
	go StartServer(datasetController)
	select {}
}
