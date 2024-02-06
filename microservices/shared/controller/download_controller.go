package controller

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/response"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

type DownloadController struct {
	app            *fiber.App
	storageService service.IStorageService
}

func NewDownloadController(storageService service.IStorageService) *DownloadController {
	return &DownloadController{
		storageService: storageService,
	}
}

func (dc *DownloadController) DownloadExport(c *fiber.Ctx) error {
	filename := c.Params("filename")
	log.WithFields(log.Fields{
		"service": "download",
	}).Info("Download started for file ", filename, " from IP: ", c.IP())
	object, err := dc.storageService.GetObject(context.Background(), configs.EnvExportBucketName(), filename)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The requested does not exist or has not been exported."))
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Set(fiber.HeaderContentType, "application/gzip")

	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to download file."))
	}
	return nil
}

func (dc *DownloadController) GetDownloadModels(c *fiber.Ctx) error {
	bucketStruc, err := dc.storageService.GetBucketStructure(configs.EnvModelBucketName(), "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to list buckets."))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(bucketStruc))
}

func (dc *DownloadController) GetDownloadExport(c *fiber.Ctx) error {
	bucketStruc, err := dc.storageService.GetBucketStructure(configs.EnvExportBucketName(), "")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to list buckets."))
	}
	return c.Status(http.StatusOK).JSON(response.NewSuccessResponse(bucketStruc))
}

func (dc *DownloadController) GetZippedModel(c *fiber.Ctx) error {
	modelName := c.Params("modelName")
	log.WithFields(log.Fields{
		"service": "DownloadController",
	}).Info("Download started for model ", modelName, " from IP: ", c.IP())
	zippedModelName := c.Params("zippedModelName")
	object, err := dc.storageService.GetObject(context.Background(), configs.EnvModelBucketName(), modelName+"/"+zippedModelName)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The requested does not exist or has not been exported."))
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", zippedModelName))
	c.Set(fiber.HeaderContentType, "application/gzip")

	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to download file."))
	}
	return nil
}

func (dc *DownloadController) GetModelData(c *fiber.Ctx) error {
	modelName := c.Params("modelName")
	trainingSet := c.Params("trainingSet")
	fileName := c.Params("filename")
	log.WithFields(log.Fields{
		"service": "DownloadController",
	}).Info("Download started for model ", modelName, ", trainingSet ", trainingSet, " and fileName ", fileName, " from IP: ", c.IP())
	object, err := dc.storageService.GetObject(context.Background(), configs.EnvModelBucketName(), modelName+"/"+trainingSet+"/"+fileName)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(response.NewErrorResponse("The requested does not exist or has not been exported."))
	}
	defer object.Close()
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	c.Set(fiber.HeaderContentType, contentType)
	if _, err = io.Copy(c.Response().BodyWriter(), object); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.NewErrorResponse("Failed to download file."))
	}
	return nil
}

func (dc *DownloadController) Start() {
	dc.app = fiber.New()
	dc.app.Use(cors.New())
	dc.app.Use(logger.New())
	dc.app.Get("/transfer/health", utils.GetHealth)
	dc.app.Get("/transfer/export/:fileName", dc.DownloadExport)
	dc.app.Get("/transfer/export", dc.GetDownloadExport)
	dc.app.Get("/transfer/models", dc.GetDownloadModels)
	dc.app.Get("/transfer/models/:modelName/:zippedModelName", dc.GetZippedModel)
	dc.app.Get("/transfer/models/:modelName/:trainingSet/:filename", dc.GetModelData)
	dc.app.Listen(":8081")
}

func (dc *DownloadController) Stop() {
	dc.app.Shutdown()
}
