package main

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/controller"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/storage"
)

func main() {
	ctx := context.Background()
	minioClient := configs.ConnectMinio()
	minioStorageAdapter := storage.NewMinIOStorageAdapter(minioClient)
	storageService := service.NewMinIOStorageService(minioStorageAdapter)
	storageService.VerifyBucketExists(ctx, configs.EnvExportBucketName())
	storageService.VerifyBucketExists(ctx, configs.EnvModelBucketName())
	downloadController := controller.NewDownloadController(storageService)
	downloadController.Start()
}
