package dataset

import (
	"context"
	"fmt"
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

func Main() {
	ctx := context.TODO()
	requestRepo := repository.NewMongoRequestRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	requestService := service.NewRequestService(requestRepo)
	datasetRepo := repository.NewMongoDatasetRepository(configs.GetDatabase(configs.ConnectDB(ctx)))
	datasetService := service.NewDatasetService(datasetRepo)

	datasetCalculationJob := job.NewDatasetMetricJob(datasetService, requestService)

	datasetController := controller.NewDatasetController(datasetService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/datasets/health", utils.GetHealth)
	app.Post("/datasets", datasetController.CreateDataset)
	app.Get("/datasets", datasetController.GetAllDatasets)

	c := cron.New()
	_, _ = c.AddFunc("@hourly", func() {
		datasetCalculationJob.Execute()
	})
	c.Start()
	go func() {
		if err := app.Listen(":8081"); err != nil {
			fmt.Println("Fiber failed to start:", err)
		}
	}()
	select {}
}
