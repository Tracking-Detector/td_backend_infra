package main

import (
	"context"
	"tds/shared/configs"
	"tds/shared/controller"
	"tds/shared/service"
	"tds/shared/storage"
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
