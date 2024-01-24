package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/service"
	"tds/shared/storage"
	"tds/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	ctx := context.TODO()
	minioClient := configs.ConnectMinio()
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	storageService := service.NewMinIOStorageService(minioStorageAdapter)
	storageService.VerifyBucketExists(ctx, configs.EnvExportBucketName())
	storageService.VerifyBucketExists(ctx, configs.EnvModelBucketName())
	downloadController := controller.NewDownloadController(storageService)
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/transfer/health", utils.GetHealth)
	app.Get("/transfer/export/:fileName", downloadController.DownloadExport)
	app.Get("/transfer/export", downloadController.GetDownloadExport)
	app.Get("/transfer/models", downloadController.GetDownloadModels)
	app.Get("/transfer/models/:modelName/:zippedModelName", downloadController.GetZippedModel)
	app.Get("/transfer/models/:modelName/:trainingSet/:filename", downloadController.GetModelData)
	app.Listen(":8081")
}
